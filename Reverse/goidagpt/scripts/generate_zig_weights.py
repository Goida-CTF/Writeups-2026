"""
Навайбкодил
"""

import json
import sys


def format_f32_array(arr, indent=4):
    indent_str = " " * indent
    if len(arr) == 0:
        return ".{}"

    elements = ", ".join(f"{x:.6f}" for x in arr)
    return f".{{ {elements} }}"


def format_2d_f32_array(arr, indent=4):
    indent_str = " " * indent
    if len(arr) == 0:
        return ".{}"

    lines = [".{"]
    for row in arr:
        row_str = format_f32_array(row, 0)
        lines.append(f"{indent_str}    {row_str},")
    lines.append(f"{indent_str}}}")
    return "\n".join(lines)


def format_u8_array(arr, indent=4):
    indent_str = " " * indent
    if len(arr) == 0:
        return ".{}"

    # Format in groups of 16 for readability
    lines = [".{"]
    for i in range(0, len(arr), 16):
        chunk = arr[i : i + 16]
        hex_strs = ", ".join(f"0x{b:02x}" for b in chunk)
        lines.append(f"{indent_str}    {hex_strs},")
    lines.append(f"{indent_str}}}")
    return "\n".join(lines)


def generate_zig_weights(weights_file="weights.json"):
    with open(weights_file, "r") as f:
        data = json.load(f)

    weights = data["weights"]
    W1 = weights["W1"]
    W2 = weights["W2"]
    W3 = weights["W3"]
    b1 = weights["b1"]
    b2 = weights["b2"]
    b3 = weights["b3"]
    encrypted_flag = data["encrypted_flag"]
    flag_length = data["flag_length"]

    zig_code = f"""// Auto-generated weight arrays
pub const IO_DIM = 8;
pub const HIDDEN_DIM = 16;

pub const W1: [HIDDEN_DIM][IO_DIM]f32 = {format_2d_f32_array(W1)};

pub const W2: [HIDDEN_DIM][HIDDEN_DIM]f32 = {format_2d_f32_array(W2)};

pub const W3: [IO_DIM][HIDDEN_DIM]f32 = {format_2d_f32_array(W3)};

pub const b1: [HIDDEN_DIM]f32 = {format_f32_array(b1)};

pub const b2: [HIDDEN_DIM]f32 = {format_f32_array(b2)};

pub const b3: [IO_DIM]f32 = {format_f32_array(b3)};

pub const encrypted_flag: [{len(encrypted_flag)}]u8 = {format_u8_array(encrypted_flag)};

pub const flag_length: usize = {flag_length};
"""

    return zig_code


def main():
    zig_code = generate_zig_weights("weights.json")

    output_file = "src/weights.zig"
    with open(output_file, "w") as f:
        f.write(zig_code)

    print(f"Zig generation done: {output_file}")


if __name__ == "__main__":
    main()
