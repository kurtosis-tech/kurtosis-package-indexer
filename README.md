Kurtosis package indexer
========================

Kurtosis package indexer is a backend services searching for Kurtosis packages in Github and storing them in memory.
Right now it is consumed by Kurtosis Frontend to power Kurtosis Packages Catalog.

Implementation details
----------------------

The service simply runs a job periodically to search for all Kurtosis Packages currently existing on Github.
- The background job runs every two hours. Results are stored in memory for now. I.e. restarting the service will re-run the job
- It searches for `kurtosis.yml` files on Github. It then checks the `kurtosis.yml` file can be parsed, and there is a valid `main.star` file next to it. Any folder not matching those criteria will be discarded
- The searches on Github needs to be authenticated. Right now the service expects a valid Github token stored as the environment variable `GITHUB_USER_TOKEN`
