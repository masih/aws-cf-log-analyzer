# aws-cf-log-analyzer

AWS CloudFront Log Analyzer that accepts two arguments:
1. Path to a directory containing `.gz` CloudFront logs, and
2. A plain text file containing the IP target addresses to search for.

It then scans each of the log files and prints the number of target IP addresses found in each of the `.gz` files along with total found.


## Usage

1. Enable AWS CloudFront logging first; see [docs](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/AccessLogs.html#access-logs-analyzing).
2. Once logs have accumulated, sync the logs to your local machine:
   ```shell
   aws s3 sync s3://<aws-cf-log-bucket>/<prefix-if-any> . 
   ```
3. Prepare a plain text file containing the target IPs. This can be in any format. The tool uses regular expression to extract IPs from the text file.
4. Install the analyzer:
   ```shell
   go install github.com/masih/aws-cf-log-analyzer@latest
   ```
5. Run the analysis:
   ```shell
   aws-cf-log-analyzer <path-to-gz-logs-directory> <path-to-target-ips-file>
   ```

You should see an output similar to the following:

```text
2023/04/24 11:38:31 found 121 IP addresses in /path/to/ips.txt
2023/04/24 11:38:31 found 1232 IPs in ABCDEFGHI.2023-04-23-20.0b9277f7.gz
2023/04/24 11:38:32 found 0 IPs in ABCDEFGHI.2023-04-23-20.114bfb34.gz
2023/04/24 11:38:32 found 0 IPs in ABCDEFGHI.2023-04-23-20.3c8527a5.gz
2023/04/24 11:38:33 found 2 IPs in ABCDEFGHI.2023-04-23-20.3da42292.gz
2023/04/24 11:38:34 found 8 IPs in ABCDEFGHI.2023-04-23-20.48d1a4df.gz
2023/04/24 11:40:19 found total of 1242 IPs across 5 .gz files
```
