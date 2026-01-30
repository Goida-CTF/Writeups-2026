use clap::Parser;
use fancy_regex::Regex;

struct Computer {
    prime_re: Regex,
    dot_product_re: Regex,
}

impl Computer {
    fn new() -> Self {
        Self {
            prime_re: Regex::new(r"^.?$|^(..+?)\1+$").unwrap(),
            dot_product_re: Regex::new(include_str!("regex.txt")).unwrap(),
        }
    }

    fn np(&self, s: &str) -> String {
        let mut s = s.to_string();
        s.push('1');

        while self.prime_re.is_match(&s).unwrap_or(false) {
            s.push('1');
        }
        s
    }

    fn ns(&self, n: u64) -> String {
        "1".repeat(n as usize)
    }

    fn c(&self, v: &str, coeffs: &str, target: &str) -> bool {
        let to_check = v.to_string() + "@" + coeffs + "=" + target;
        self.dot_product_re.is_match(&to_check).unwrap_or(false)
    }
}

#[derive(Parser)]
struct Cli {
    flag: String,
}

fn main() {
    let coeffs = include!("coeffs.rs");
    let computer = Computer::new();
    let cli = Cli::parse();

    let mut flag = cli.flag;

    if flag.len() != 42 {
        println!("nep");
        return;
    }

    while flag.len() % 4 != 0 {
        flag.push('\x00');
    }

    let valid = flag
        .as_bytes()
        .chunks(4)
        .map(|v| {
            "[".to_owned()
                + &v.iter()
                    .map(|c| computer.ns(*c as u64))
                    .collect::<Vec<String>>()
                    .join(",")
                + "]"
        })
        .zip(coeffs)
        .all(|(v, coeffs)| {
            coeffs.iter().all(|(c, t)| {
                computer.c(
                    &v,
                    &("[".to_owned()
                        + &c.iter()
                            .map(|c| computer.np(&computer.ns(*c as u64)))
                            .collect::<Vec<String>>()
                            .join(",")
                        + "]"),
                    t,
                )
            })
        });
    if valid {
        println!("GOIDA BRATYA!!!!!!!");
    } else {
        println!("nea");
    }
}
