# echocat's mongod vfox plugin

This provides a vfox plugin to use `mongod` and `mongos` inside your project as tools. It is compatible with [vfox](https://vfox.dev/) itself and [MISE](https://mise.jdx.dev).

## TOC

1. [Environment](#environment)
2. [Usage](#usage) ...with [vfox](#with-vfox) or [MISE](#with-mise)
3. [FAQ](#faq)
4. [Contributing](#contributing)
5. [License](#license)

## Environment

The actual target and architecture for the given MongoDB version will be resolved based on your operating-system **and** architecture **and** distribution (if applicable) **and** os/distribution version (if applicable). This is especially important for Linux distributions, where a big variety of packages for different distributions exists for `mongod`/`mongos`.

But: If there is something missing, please [create a ticket](https://github.com/echocat/vfox-mongod/issues) that we can implement it.

### Overrides

| Environment Variable | Description |
| -- | -- |
| `MONGOD_TARGET` | This will override the automatically resolved target (like `ubuntu2404`, `windows`, `macos`, ...). All available lists of targets are listed inside [downloads.mongodb.org/full.json](https://downloads.mongodb.org/full.json). |
| `MONGOD_ARCH` | This will override the architecture. Can be (currently) `amd64`, `arm`, `arm64`, `s390x` and `ppc64le`. |

## Usage

1. [With vfox](#with-vfox)
2. [With MISE](#with-mise)

### With vfox

1. Install plugin
   ```shell
   vfox add --source https://github.com/echocat/vfox-mongod/archive/refs/tags/<release>.zip mongod
   ```

2. Install the tool
    * Always latest version:
      ```shell
      vfox install mongod@latest
      ```
    * Specific version:
      ```shell
      vfox install mongod@8.2.1
      ```

3. Check if it is working
   ```shell
   vfox use mongod@8.2.1
   mongod --version
   ```

### With MISE

You have two options, either use it as a [vfox plugin / directly as a MISE tool](#vfox-plugin) or as an [MISE plugin](#mise-plugin).

#### vfox Plugin

1. Add a handy alias
   ```shell
   mise alias set vfox:echocat/vfox-mongod mongod
   ```

2. Install the tool into your project
    * Always latest version:
      ```shell
      mise use mongod@latest
      ```
    * Specific version:
      ```shell
      mise use mongod@8.2.1
      ```

3. Check if it is working
   ```shell
   mise exec -- mongod --version
   ```

#### MISE Plugin

1. Install the plugin (ℹ️ this only needs to be done once per user - **not** per project - it will be installed inside your current user)
   ```shell
   mise plugin install mongod https://github.com/echocat/vfox-mongod
   ```

2. Install the tool into your project
    * Always latest version:
      ```shell
      mise use mongod@latest
      ```
    * Specific version:
      ```shell
      mise use mongod@8.2.1
      ```

3. Check if it is working
   ```shell
   mise exec -- mongod --version
   ```

## FAQ

### Can I override the target which will be picked?

See [Environment -> Overrides](#overrides).

### What is vfox?

See [vfox.dev](https://vfox.dev)

### What is MISE?

See [mise.jdx.dev](https://mise.jdx.dev)

### Where does this plugin know which versions of mongod are available?

It uses the [downloads.mongodb.org/full.json](https://downloads.mongodb.org/full.json) as source which is maintained by MongoDB itself.

## Contributing

**vfox-mongod** is an open source project by [echocat](https://echocat.org). So if you want to make this project even better, you can contribute to this project on [Github](https://github.com/echocat/vfox-mongod) by [fork us](https://github.com/echocat/vfox-mongod/fork).

If you commit code to this project, you have to accept that this code will be released under the [license](#license) of this project.

## License

See the [LICENSE](LICENSE) file.
