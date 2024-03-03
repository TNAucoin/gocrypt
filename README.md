# gocrypt

Small CLI tool written in go to encrypt files using AWS KMS.
Useful for incorporating encrypted files into a CI/CD pipeline, or a local
development setup. Supports both direct and role separation usage.

> Currently only supports symmetric encryption.

## Usage

Default Behaviors:

- If you don't specify the output file, the encrypted file will be saved in the same directory as the input file.
- If you don't specify the role, the tool will use the default profile role.
- If a file already exists in the output directory with the same name as the input file, the tool will rename the output file with the extension specified in the `--ext` flag.
- If a encrypted file already exists in the output directory with the same name as the new output, the tool will rename that file with `.old` then create the new encrypted file. (e.g. `filename.something.encrpyted.old`)

### Encrypt

#### Flags
| Flag        | Flag Short                   | Description | Required | Default Value |
|-------------|------------------------------| --- | --- |---------------|
| `--key`     | `-k` | KMS Key Alias                | X | " "
| `--output`  | `-o` | Output file path             | | " "
| `--path`    | `-P`   | Path to the file to encrypt  | X |" "
| `--profile` | `-p`   | AWS Profile name             | X | default       
| `--region`  | `-r`    | AWS Region                   | | us-east-1     
| `--role`     | `-R`    | Role to assume               | | " "
| `--ext` | `-e` | File extension to encrypt    | | "encrypted"
Example:
```bash
gocrypt encrypt -k alias/gocrypt -o ./output/path/newFilename.something -P input/path/filename.something -p awsProfileName -r aws-region -R arn:aws:iam::RoleToAssume -e encoded
```

Encrypted file:
`./output/path/newFilename.something.encoded`

> If your AWS profile has KMS permissions there is no need to supply the `--role` flag. If 
> your AWS profile does not have KMS permissions, you must supply the `--role` flag to a role
> with the correct permissions KMS permissions (kms::encrypt)

### Decrypt

#### Flags
| Flag        | Flag Short                   | Description | Required | Default Value |
|-------------|------------------------------| --- | --- |---------------|
| `--output`  | `-o` | Output file path             | | " "
| `--path`    | `-P`   | Path to the file to encrypt  | X |" "
| `--profile` | `-p`   | AWS Profile name             | X | default       
| `--region`  | `-r`    | AWS Region                   | | us-east-1     
| `--role`     | `-R`    | Role to assume               | | " "
| `--ext` | `-e` | File extension to encrypt    | | "encrypted"

Example:

```bash
gocrypt decrypt -k alias/gocrypt -o ./output/path/newFilename.something -P input/path/filename.something.
```

