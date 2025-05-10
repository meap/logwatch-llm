# Logwatch LLM reports

## Installation from Source

> **Requirements:** Go 1.24 or newer must be installed on your system. See [https://go.dev/doc/install](https://go.dev/doc/install) for details.

To install Logwatch LLM reports from source on Linux, follow these steps:

1. **Clone the Repository**
   
   ```sh
   git clone https://github.com/meap/logwatch-llm.git
   cd logwatch-llm
   ```

2. **Build the Project**
   
   Use the provided Makefile to build the binary:
   
   ```sh
   make build
   ```
   
   The compiled binary will be located in the `bin/` directory as `logwatch-llm`.

3. **(Optional) Install the Binary System-wide**
   
   To make the CLI available system-wide:
   
   ```sh
   sudo cp bin/logwatch-llm /usr/local/bin/
   ```

### Permissions and Ownership After Install

After copying the binary to `/usr/local/bin`, ensure it is executable and securely owned:

```sh
sudo chmod 755 /usr/local/bin/logwatch-llm
sudo chown root:root /usr/local/bin/logwatch-llm
```

This makes the binary usable by all users and prevents unauthorized modifications.

## Running with Logwatch via Cron

To automate Logwatch analysis and reporting, you can run Logwatch and pipe its output directly to Logwatch LLM using a cron job. For example, to run this daily and email the report:

1. Edit your crontab:

   ```sh
   crontab -e
   ```

2. Add a line like this (adjust paths and email as needed):

   ```sh
   0 7 * * * /usr/sbin/logwatch | /usr/local/bin/logwatch-llm --model gpt-4o --email your@email.com
   ```

This will run Logwatch at 7:00 AM every day, analyze the output with the LLM, and send the report to your email.
