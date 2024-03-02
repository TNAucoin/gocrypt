# gocrypt

Small CLI tool written in go to encrypt files using AWS KMS.
Useful for incorporating encrypted files into a CI/CD pipeline, or a local
development setup.

Basic usage:
```bash
gocrypt encrypt -k alias/gocrypt -o ./output/path/filename.something -P input/path/filename.something -p awsProfileName -r aws-region -R arn:aws:iam::RoleToAssume
```
> If you don't specify a role then the tool will assume the role of the provided AWS Profile

