# Changelog

## [0.0.5](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.4...0.0.5) (2023-08-30)


### Features

* Add endpoint to manually trigger a reindex ([#20](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/20)) ([d6c777a](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/d6c777a573329ee16389fb007ed1880b7e6ec1ec))
* Add package URL in returned Kurtosis package objects ([#17](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/17)) ([d1f6e6b](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/d1f6e6b294e93b3792bd8a69df4c28df9dc93c14))
* Parse package description from kurtosis.yml ([#21](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/21)) ([65d6e81](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/65d6e811f4b7f62c34ef8dcf449992638b563bd5))

## [0.0.4](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.3...0.0.4) (2023-08-29)


### Bug Fixes

* Package with a `kurtosis.yml` file with an empty name are not considered invalid ([#15](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/15)) ([6b989f4](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/6b989f4c6832012603ae4686e92f8c1524b840f1))

## [0.0.3](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.2...0.0.3) (2023-08-29)


### Features

* Improve Kurtosis packages search ([#10](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/10)) ([a1818ca](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/a1818ca1cbd990e845ae7bcd79af0a3dc58f83be))

## [0.0.2](https://github.com/kurtosis-tech/kurtosis-package-indexer/compare/0.0.1...0.0.2) (2023-08-28)


### Features

* Add connect-go for the FE to be able to talk to the indexer ([#5](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/5)) ([d796aa3](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/d796aa39c4b1d1d3c00c5ed2ab6defd4338245d5))
* Add GetPackages endpoint ([#1](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/1)) ([9b9f63a](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/9b9f63a9ff7bbd99289650e2a962ad9163fdf102))
* Add Github auth based on user token ([#7](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/7)) ([b2d4c4a](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/b2d4c4a762a35e1faf46acc440677ffcdbcb7882))
* Crawl Github to fetch existing Kurtosis packages ([#4](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/4)) ([bc8d573](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/bc8d573036cfbb47f86120132a49c2980431f921))
* Return Github repo stars as well ([#8](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/8)) ([c9fbe5d](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/c9fbe5d22f19457b06c5971c66441416f123467b))
* The crawler now parses the main.star file ([#6](https://github.com/kurtosis-tech/kurtosis-package-indexer/issues/6)) ([3f09acb](https://github.com/kurtosis-tech/kurtosis-package-indexer/commit/3f09acbf1cdab73dcbe7ba8ab987d327abdc68bd))
