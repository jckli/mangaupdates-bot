# Changelog

## [3.2.0](https://github.com/jckli/mangaupdates-bot/compare/v3.1.0...v3.2.0) (2026-08-04)


### Features

* identify manga search results ([792bf63](https://github.com/jckli/mangaupdates-bot/commit/792bf63e225e99d00c1e24a8e03b6d4a6118e2fe))
* unify server management permissions ([35f6bee](https://github.com/jckli/mangaupdates-bot/commit/35f6beebeeb57c7bf1ab64a19306b4ac40ae977d))
* use unified manga metadata ([f0b73c6](https://github.com/jckli/mangaupdates-bot/commit/f0b73c66aeb7c7c1df726855c67a2e2b58a55d9b))


### Bug Fixes

* handle missing manga search years ([9819e65](https://github.com/jckli/mangaupdates-bot/commit/9819e65eff33a18de1a6e8d6064ec969c5abf6a1))
* refine manga result presentation ([bef9904](https://github.com/jckli/mangaupdates-bot/commit/bef9904351ce92e049765dc1af606e153f1de0bb))
* route manga search selections ([c73a6ad](https://github.com/jckli/mangaupdates-bot/commit/c73a6ad49fecb8fa29dbb9e7edf3a941ef5182be))

## [3.1.0](https://github.com/jckli/mangaupdates-bot/compare/v3.0.0...v3.1.0) (2026-08-01)


### Features

* add internal /internal/status API for metrics ([71b6abd](https://github.com/jckli/mangaupdates-bot/commit/71b6abdf0242f6cc1676705871ab1d183ed18216))
* add ping role option to server slash command and bridge notification payload ([63071a9](https://github.com/jckli/mangaupdates-bot/commit/63071a9b53c2e51b6918fe204dbac30c9f852712))
* auto guild delete when bot gets kicked ([3663d65](https://github.com/jckli/mangaupdates-bot/commit/3663d652c10d4d422a6fc3d5786ee71973a60647))
* include start_time in status payload for client side uptime tracking ([d9149c5](https://github.com/jckli/mangaupdates-bot/commit/d9149c595410e80449a6832cef20ffbe1ec17dc8))


### Bug Fixes

* auto publish the all chapter updates ([b14ac5c](https://github.com/jckli/mangaupdates-bot/commit/b14ac5ccbccdc4a5bc65d236903f2f0557ca868c))
* correct disgo shard ready status enum, start stats worker immediately, and include total users ([9c80cd1](https://github.com/jckli/mangaupdates-bot/commit/9c80cd196c26d33b47116463e17b98fe7d607db9))
* handle guild update member count ([8704b61](https://github.com/jckli/mangaupdates-bot/commit/8704b61471b1e6db310a1e56c3206b0d5e3576fc))
* improve non-admin role confirmation text in server role commands ([519109c](https://github.com/jckli/mangaupdates-bot/commit/519109c09d826dbb263b9dac7a1d2fdf841cdf90))
* **interaction:** run component selection asynchronously to avoid 3s Discord timeout ([195197b](https://github.com/jckli/mangaupdates-bot/commit/195197bb9311fac6a5862478a65bfe40529bbc30))
* **interaction:** send error embed on deferred failure and add config caching ([cdc85a9](https://github.com/jckli/mangaupdates-bot/commit/cdc85a9028f1705ccf671bb5d1e0d1e949bdf322))
* reorder slash command options for server role set to type first ([f8e2a80](https://github.com/jckli/mangaupdates-bot/commit/f8e2a8065feea529d85c8fef80896cf38745e4f9))
* retrieve gateway shard latency in ping command ([50f72ff](https://github.com/jckli/mangaupdates-bot/commit/50f72ffded95f9c8990db4d0ba7501ff594b932e))
* set fetch-depth 0 to correctly resolve git tags in docker build ([ba24043](https://github.com/jckli/mangaupdates-bot/commit/ba240434fac210e084d5b14ce2502ae7d98228a3))
