Kurtosis package indexer
========================

Kurtosis package indexer is a backend services searching for Kurtosis packages in Github and storing them in memory.
Right now it is consumed by Kurtosis Frontend to power Kurtosis Packages Catalog.

Implementation details
----------------------

The service simply runs a job periodically to search for all Kurtosis Packages currently existing on Github.
- The background job runs every two hours. Results are stored in memory for now. I.e. restarting the service will re-run the job
- It searches for `kurtosis.yml` files on Github. It then checks the `kurtosis.yml` file can be parsed, and there is a valid `main.star` file next to it. Any folder not matching those criteria will be discarded

Github authentication
---------------------

The searches ran on Github need to be authenticated. There're two ways Kurtosis Package Indexer will authenticate itself
on Github.
Right now, the indexer first tries reading the `GITHUB_USER_TOKEN` environment variable and if it's empty, it falls back
to the S3 bucket option.

### Using a Github token via an environment variable
This is the simplest. The indexer expects a valid Github token stored inside the environment variable `GITHUB_USER_TOKEN`

### Using a file stored inside an S3 bucket
The indexer can also get the Github token from a file stored inside an S3 bucket.
The file storing the Github token should be named `github-user-token.txt` and it should contain only the Github token 
on as plain text.

To access this file, the indexer will require the following environment variables to be set:
- `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` for AWS authentication
- `AWS_BUCKET_REGION` and `AWS_BUCKET_NAME` to identify the AWS S3 bucket. The user linked to the key above needs to  
be able to do `GetObject` on this bucket
- `AWS_BUCKET_FOLDER` (optional) in case the file `github-user-token.txt` is located inside a folder in this S3 bucket

Data persistence
----------------

The Kurtosis packages information are stored by default in-memory. Everytime the indexer is restarted, it re-runs the
Github searches to fetch the latest information about the packages on Github

There's also the option of persisting the data to a [bolt](https://github.com/etcd-io/bbolt) key value store, so that 
services can be restarted keeping the data intact. To use it, the environment variable `BOLT_DATABASE_FILE_PATH` can 
be set to point to a file on disk that bolt will use to store the data. If the indexer is being run in a container, a 
persistent volume should be used to fully benefit from this feature.
