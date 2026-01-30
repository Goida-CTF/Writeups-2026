#!/usr/bin/env python3
"""
Generate a pure regex that validates unary dot products:
    [a1,a2,...,an]@[b1,b2,...,bn]=result

where result = a1*b1 + a2*b2 + ... + an*bn (all in unary: 5 = "11111")

Uses the "extremely slow" backtracking approach where the regex engine
tries all possible decompositions of the result.

The regex validates:
1. Result is composed of a_i values only  
2. For each product i, the number of a_i copies equals number of 1s in b_i
3. Products appear in order (non-interleaved)

Uses the `regex` module's captures counting for exact verification.
"""

import regex
import sys


def generate_dot_product_regex_space_optimized(n: int = 4, max_b_len: int = 10) -> str:
    """
    Space-optimized regex for n=4 dot product.
    
    Uses full enumeration (which is required for correctness) but with
    optimizations to reduce pattern size:
    
    1. Uses compact b patterns (just '1' * k instead of complex structures)
    2. Uses efficient r patterns with exact counts
    3. Skips redundant lookahead prefix matching
    4. Smaller default max_b_len (10 instead of 50)
    
    For max_b_len=10: 11^4 = 14,641 combinations, ~100KB pattern
    For max_b_len=5:  6^4 = 1,296 combinations, ~10KB pattern
    
    Args:
        n: Number of dimensions (must be 4)
        max_b_len: Maximum length of each b_i value (default 10)
    
    Returns:
        A pure regex pattern that validates the dot product exactly
    """
    if n != 4:
        raise ValueError("Space-optimized version currently only supports n=4")
    
    from itertools import product as cart_product
    
    # Capture all a values
    a_caps = ','.join([f'(?P<a{i}>1*)' for i in range(1, 5)])
    
    # Build alternation for all (k1, k2, k3, k4) combinations
    # Each alternative specifies exact b pattern and exact r pattern
    alts = []
    for k1, k2, k3, k4 in cart_product(range(max_b_len + 1), repeat=4):
        # b pattern: exact counts
        b_pat = ','.join(['1' * k if k > 0 else '' for k in (k1, k2, k3, k4)])
        
        # r pattern: exact repetitions using quantifiers
        r_parts = []
        for i, k in enumerate((k1, k2, k3, k4), 1):
            if k > 0:
                r_parts.append(f'(?P=a{i}){{{k}}}')
        r_pat = ''.join(r_parts)
        
        # Use lookahead to verify before matching
        # This ensures we match the right branch
        alts.append(f'{b_pat}\\]={r_pat}$')
    
    b_section = '(?:' + '|'.join(alts) + ')'
    
    return f'^\\[{a_caps}\\]@\\[{b_section}'


def generate_dot_product_regex(n: int, max_b_len: int = 50) -> str:
    """
    Generate a SINGLE pure regex for n-dimensional dot product validation.
    
    Uses alternation to enumerate all possible b-length combinations.
    The pattern verifies [a1,...,an]@[b1,...,bn]=result
    where result = (a1)^{len(b1)} ++ (a2)^{len(b2)} ++ ... ++ (an)^{len(bn)}
    
    This is the "extremely slow" backtracking approach because the regex
    engine tries all possible decompositions via alternation.
    
    Args:
        n: Number of dimensions
        max_b_len: Maximum length of each b_i value to support
    
    Returns:
        A pure regex pattern that validates the dot product exactly
    """
    
    if n < 1:
        raise ValueError("n must be at least 1")
    
    # Capture a values
    a_captures = ','.join([f'(?P<a{i}>1*)' for i in range(1, n + 1)])
    
    # Build alternation for all b-length combinations
    # For n=1: b can be "", "1", "11", ... -> result is "", a, aa, ...
    # For n>1: b is (b1, b2, ...) -> result is r1 ++ r2 ++ ...
    
    def build_b_alts(dim: int) -> list:
        """Build list of (b_pattern, r_pattern) for dimension dim."""
        alts = []
        for k in range(0, max_b_len + 1):
            b_pat = '1' * k if k > 0 else ''
            r_pat = f'(?P=a{dim}){{{k}}}' if k > 0 else ''
            alts.append((b_pat, r_pat))
        return alts
    
    if n == 1:
        # For n=1, build simple alternation
        alts = []
        for k in range(0, max_b_len + 1):
            b_pat = '1' * k if k > 0 else ''
            r_pat = f'(?P=a1){{{k}}}' if k > 0 else ''
            alts.append(f'(?={b_pat}\\]={r_pat}$){b_pat}')
        
        b_section = '(?:' + '|'.join(alts) + ')'
        return f'^\\[(?P<a1>1*)\\]@\\[{b_section}\\]=(?:(?P=a1))*$'
    
    else:
        # For n>1, we need to handle each product separately
        # Each b_i can be any length, and corresponding r_i must have that many copies
        
        # Build b section with named captures
        b_captures = ','.join([f'(?P<b{i}>1*)' for i in range(1, n + 1)])
        
        # Build result verification using lookahead for each product
        # The trick: use lookahead to verify the structure before matching
        
        # For ordered products: result = r1 ++ r2 ++ ... ++ rn
        # Each r_i = (a_i)^{len(b_i)}
        
        # Use lookahead with alternation for each product's count
        lookahead_parts = []
        for i in range(1, n + 1):
            # Previous products: match greedily
            prev = ''.join([f'(?:(?P=a{j}))*' for j in range(1, i)])
            # This product: alternation for exact count
            alts = []
            for k in range(0, max_b_len + 1):
                b_check = '1' * k if k > 0 else ''
                r_check = f'(?P=a{i}){{{k}}}' if k > 0 else ''
                # Lookahead checks: b_i has exactly k ones AND r_i has exactly k copies
                alts.append(
                    f'(?=\\[.*\\]@\\[' + 
                    ','.join(['1*'] * (i-1) + [b_check] + ['1*'] * (n-i)) + 
                    f'\\]={prev}{r_check})'
                )
            # Next products: any
            next_prods = ''.join([f'(?:(?P=a{j}))*' for j in range(i + 1, n + 1)])
            lookahead_parts.append(f'(?:{"|".join(alts)})')
        
        # Combine into full pattern
        # This is getting complex, use simpler structure
        
        # Simpler approach: enumerate all combinations of b lengths
        # For small n and max_b_len, this is feasible
        
        from itertools import product as cart_product
        
        all_alts = []
        for lengths in cart_product(range(max_b_len + 1), repeat=n):
            # lengths = (k1, k2, ..., kn) where ki = len(bi)
            b_pat = ','.join(['1' * k if k > 0 else '' for k in lengths])
            r_pat = ''.join([f'(?P=a{i}){{{lengths[i-1]}}}' if lengths[i-1] > 0 else '' 
                            for i in range(1, n + 1)])
            all_alts.append(f'(?={b_pat}\\]={r_pat}$){b_pat}')
        
        b_section = '(?:' + '|'.join(all_alts) + ')'
        r_section = ''.join([f'(?:(?P=a{i}))*' for i in range(1, n + 1)])
        
        return f'^\\[{a_captures}\\]@\\[{b_section}\\]={r_section}$'


def match_dot_product(pattern: str, text: str, n: int) -> bool:
    """
    Match text against dot product pattern and verify using backtracking.
    
    This implements the "extremely slow" approach where we:
    1. Parse the structure with regex
    2. Use backtracking to verify result = a1*b1 + a2*b2 + ... + an*bn
    
    The backtracking tries to decompose the result into ordered products,
    each being a_i repeated len(b_i) times.
    
    Returns True if:
    1. Text matches the structure [a1,...,an]@[b1,...,bn]=result
    2. Result can be decomposed as (a1)^{len(b1)} ++ (a2)^{len(b2)} ++ ...
    """
    m = regex.match(pattern, text)
    if not m:
        return False
    
    # Extract values
    a_vals = [m.group(f'a{i}') for i in range(1, n + 1)]
    b_vals = [m.group(f'b{i}') for i in range(1, n + 1)]
    result = m.group('result')
    
    # Backtracking verification: try to decompose result into products
    def try_decompose(product_idx: int, remaining: str) -> bool:
        """Try to decompose remaining result starting at product_idx."""
        if product_idx > n:
            return remaining == ''  # All products allocated, result exhausted
        
        a = a_vals[product_idx - 1]
        b = b_vals[product_idx - 1]
        expected_copies = len(b)
        
        if len(a) == 0:
            # a is empty, product contributes nothing regardless of b
            return try_decompose(product_idx + 1, remaining)
        
        # This product should contribute exactly expected_copies of a
        expected_len = expected_copies * len(a)
        if len(remaining) < expected_len:
            return False
        
        # Verify the portion matches a repeated expected_copies times
        portion = remaining[:expected_len]
        if portion != a * expected_copies:
            return False
        
        # Recurse with remaining result
        return try_decompose(product_idx + 1, remaining[expected_len:])
    
    return try_decompose(1, result)


def test_regex(pattern: str, test_cases: list, n: int) -> None:
    """Test the PURE regex pattern against a list of (input, expected) pairs."""
    print(f"Testing pattern (length={len(pattern)}):\n")
    
    passed = 0
    failed = 0
    for test_input, expected in test_cases:
        # Test the PURE regex without helper verification
        result = bool(regex.match(pattern, test_input))
        status = "✓" if result == expected else "✗"
        if result == expected:
            passed += 1
        else:
            failed += 1
        print(f"  {status} '{test_input}' -> {result} (expected {expected})")
    
    print(f"\nPassed: {passed}/{passed+failed}")


def main():
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} <n> [test] [max_b_len] [--optimized]")
        print("\nGenerate PURE regex for n-dimensional dot product validation.")
        print("Add 'test' to run test cases.")
        print("max_b_len (default 10) controls the maximum b_i length supported.")
        print("--optimized uses compact pattern for n=4 (smaller default max_b)")
        sys.exit(1)
    
    n = int(sys.argv[1])
    run_tests = 'test' in sys.argv
    use_optimized = '--optimized' in sys.argv or '-o' in sys.argv
    max_b_len = None  # Will set based on mode
    for arg in sys.argv[2:]:
        if arg.isdigit():
            max_b_len = int(arg)
    
    if use_optimized and n == 4:
        if max_b_len is None:
            max_b_len = 10  # Smaller default for optimized
        pattern = generate_dot_product_regex_space_optimized(n, max_b_len=max_b_len)
        num_alts = (max_b_len + 1) ** 4
        print(f"# Optimized n=4: {num_alts} alternations, ~{len(pattern)//1024}KB", file=sys.stderr)
    else:
        if max_b_len is None:
            max_b_len = 20  # Default for general case
        pattern = generate_dot_product_regex(n, max_b_len=max_b_len)
        print(f"# Full enumeration: {(max_b_len + 1) ** n} alternations", file=sys.stderr)
    
    print(pattern)
    
    if run_tests:
        print(f"\n{'='*60}")
        print(f"Testing regex for n={n}")
        print('='*60)
        
        if n == 1:
            test_cases = [
                ('[1]@[1]=1', True),            # 1*1 = 1
                ('[11]@[11]=1111', True),       # 2*2 = 4
                ('[111]@[11]=111111', True),    # 3*2 = 6
                ('[11]@[111]=111111', True),    # 2*3 = 6
                ('[1]@[11111]=11111', True),    # 1*5 = 5
                ('[11111]@[1]=11111', True),    # 5*1 = 5
                ('[11]@[11]=111', False),       # 2*2 ≠ 3
                ('[]@[]=', True),               # 0*0 = 0
            ]
        elif n == 2:
            test_cases = [
                ('[1,1]@[1,1]=11', True),           # 1*1 + 1*1 = 2
                ('[11,1]@[1,11]=1111', True),       # 2*1 + 1*2 = 4
                ('[11,11]@[11,11]=11111111', True), # 2*2 + 2*2 = 8
                ('[1,1]@[1,1]=111', False),         # 1+1 ≠ 3
            ]
        elif n == 3:
            test_cases = [
                ('[1,11,1]@[11,11,1111]=1111111111', True),   # 1*2 + 2*2 + 1*4 = 10
                ('[1,1,1]@[1,1,1]=111', True),                # 1+1+1 = 3
                ('[11,11,11]@[11,11,11]=111111111111', True), # 4+4+4 = 12
                ('[1,11,1]@[11,11,1111]=11111111111', False), # Wrong sum
            ]
        else:
            test_cases = []
            print("No predefined test cases for this n. Add your own!")
        
        if test_cases:
            test_regex(pattern, test_cases, n)


if __name__ == '__main__':
    main()
