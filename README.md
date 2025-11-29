# Spartan

> A simple and secure web server for SPA (single page application) or similar
> apps with static assets.

## Usage

```
spartan [flags]

Flags:
  -p, --port uint16                 The local port to listen on for incoming requests. (default 8080)
  -r, --server-path-root string     The absolute path on the server where the static content is exposed.
      --config string               The config file to use.
      --log_dir string              If non-empty, write log files in this directory (no effect when -logtostderr=true)
      --logtostderr                 log to standard error instead of files (default true)
  -d, --static-content-dir string   The path to the directory holding the static content to serve. (default "/content")
  -v, --v Level                     number for the log level verbosity
  -h, --help                        help for spartan
```

For example, the following command starts the `spartan` web server for static
content from the directory `/path/to/web` and exposing it on the sub-path `/web/`
through the HTTP interface:

```bash
spartan -d /path/to/web -r /web/
```

That is, when a browser requests for example _http://localhost:8080/web_, the
browser will automatically render the file _/path/to/web/index.html_ served by
`spartan`.

### Container images

When using the official spartan container images, `spartan` is located at
`/srv/spartan` and it expects the `config.yaml` file in the same directory, i.e.,
it tries to load `/srv/config.yaml`. However, you can of course put the YAML
configuration file in any location of your liking and instruct `spartan` to load
it from there. For example, with _docker_:

```bash
docker container run --rm -it ghcr.io/rokeller/spartan:0 --config=/path/to/my-config.yaml
```

Please note that of course you should **not** put the configuration file in the
same folder (or a subfolder) as the static content, since otherwise your
configuration file can be exposed by `spartan` itself.

You can run with the default configuration, but unless you add some content, you
won't be able see anything in a browser. The easiest way to try this is to
create a file called `index.html` in the current working directory and run
something like the following:

```bash
docker container run --rm -it \
  --mount type=bind,source=$PWD/index.html,target=/content/index.html,readonly \
  -p 8080:8080 \
  ghcr.io/rokeller/spartan:0 -v2
```

This will serve whatever you put into your `index.html` for whatever path you're
requesting on [localhost:8080](http://localhost:8080), because `spartan` treats
your `index.html` as the fallback resource for any paths that have no matching
resource in the container - the behavior desired for some SPAs. To disable this
fallback, set the `server.fallbackToIndex` property to `false`.

The additional `-v2` switch turns on logging of verbosity levels up to 2, thus
including logging for requests.

#### Image flavors

We publish two different flavors of the container image: _minimal_ (with _no_
image tag suffix, i.e., the default for each version) and _tools_ (with the
`-tools` suffix after the version). In production, you'll typically want to run
the _minimal_ flavor, as it has the minimal footprint. The _tools_ flavor offers
additional commands on the `spartan` executable that can be helpful to inspect
content served by `spartan` or headers added to responses. Those commands are
discussed in the sections below. You can also find out more by running for
example:

```bash
docker container run --rm -it ghcr.io/rokeller/spartan:0-tools -h
```

#### Command `spartan ls`

The `ls` subcommand of `spartan` lists content in the configured static content
directory of the image, which by default is `/content`. For example, if you
create your own image with the `spartan:version-tools` image as the base image,
you can inspect all content files by running

```bash
docker container run --rm -it my-image-with-spartan ls
```

#### Command `spartan headers`

The `headers` subcommand of `spartan` prints the values that `spartan` adds for
various headers in responses, as configured in the YAML configuration. This can
be useful to debug these headers without the need to run the server. For example,
if you quickly want to validate your `config.yaml` and see the various response
headers it produces, you can run

```bash
docker container run --rm -it \
    -v /path/to/my/config.yaml:/srv/config.yaml:ro \
    ghcr.io/rokeller/spartan:0-tools headers
```

The output then shows the response headers exactly as they would be presented
to a browser with the specific configuration. This can be helpful also in
combination with other tools like for example
[CSP Evaluator](https://csp-evaluator.withgoogle.com/) by Google.

#### Health endpoints

The `spartan` server offers two health endpoints.

* `GET /_spartan/live` serves as a liveness check. It always responds with HTTP
  200 and the JSON body `{"status":"ok"}` when the server is running.
* `GET /_spartan/runtime` serves some runtime information about the server, such
  as heap and stack usage, number of goroutines and garbage collections.

## Configuration

By default `spartan` is loading additional configuration from a file called
`config.yaml` in the current working directory. You can change that by using the
`--config` flag as shown above.

The config file itself has the following structure, using sample values:

```yaml
server:
  port: 8080
  staticContentDir: /content
  pathRoot: /my-spa/
  fallbackToIndex: true # Defaults to true when omitted, set to false to disable

  cache:
    defaultPolicy:
      # The cache policy to use, see below.
    routes:
      # Allows to configure different cache policies for different endpoints. See below

  security:
    contentTypeOptionsNoSniff: true
    referrerPolicy: strict-origin-when-cross-origin
    contentSecurityPolicy:
      # The content security policy to use, see below.
    permissionsPolicy:
      # The permissions policy to use, see below.
    reportingEndpoints:
      # The reporting endpoints to use for CSP reporting, see below.
    strictTransportSecurityPolicy:
      # The strict transport security policy to use, see below.
```

The `server.security` configuration section drives the security-related headers
in server responses.

The property `contentTypeOptionsNoSniff` controls the presence of the
`X-Content-Type-Options` header. When it is set to true (which is the default),
`spartan` will issue the response header `X-Content-Type-Options: nosniff`.
See also [X-Content-Type-Options header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/X-Content-Type-Options)
for more information.

The property `referrerPolicy` controls the value of the `Referrer-Policy` header.
It defaults to the value `same-origin` and it cannot be turned off. But you can
overwrite it with a referrer policy that is more to your liking as needed, by
using any of the values allowed per [Referrer-Policy header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Referrer-Policy).
Please note that not all referrer policy values are supported by all browsers.

### Cache policy

Here's how the cache policy can be configured. The block must be under
`server.cache` in the YAML file.

```yaml
defaultPolicy:
  immutable: boolean       # Set to add the 'immutable' directive to the 'cache-control' response header
  mustRevalidate: boolean  # Set to add the 'must-revalidate' directive to the 'cache-control' response header
  mustUnderstand: boolean  # Set to add the 'must-understand' directive to the 'cache-control' response header
  noCache: boolean         # Set to add the 'no-cache' directive to the 'cache-control' response header
  noStore: boolean         # Set to add the 'no-store' directive to the 'cache-control' response header
  noTransform: boolean     # Set to add the 'no-transform' directive to the 'cache-control' response header
  private: boolean         # Set to add the 'private' directive to the 'cache-control' response header
  proxyRevalidate: boolean # Set to add the 'proxy-revalidate' directive to the 'cache-control' response header
  public: boolean          # Set to add the 'public' directive to the 'cache-control' response header

  # See the 'Duration values' section below to learn how to configure these values.
  maxAge: duration               # Set the value for the 'max-age' directive.
  sharedMaxAge: duration         # Set the value for the 's-maxage' directive.
  staleIfError: duration         # Set the value for the 'stale-if-error' directive.
  staleWhileRevalidate: duration # Set the value for the 'stale-while-revalidate' directive.
```

For example, the following cache policy ...

```yaml
server:
  cache:
    defaultPolicy:
      maxAge: 168h # 1 week
      sharedMaxAge: 168h # 1 week
      mustRevalidate: true
      public: true
      staleWhileRevalidate: 24h
      staleIfError: 24h
```

... will produce a `cache-control` header like

`Cache-Control: max-age=604800, s-maxage=604800, stale-if-error=86400, stale-while-revalidate=86400, must-revalidate, public`

#### Route-matched cache policies

`spartan` also allows you configure different cache policies for different
resources. For example, some SPA bundlers generate assets (scripts, stylesheets,
fonts, images, etc.) using hash-based file names, implying that an asset's file
name never changes unless the asset input changes. This allows for virtually
infinite caching of such assets.

Currently, you can configure only path prefix-based cache policies, aside from
the above-mentioned default cache policy. For example, to configure a different
cache policy for everything in the `assets/` folder of your static content,
you could do something like the following:

```yaml
server:
  cache:
    defaultPolicy:
      # Put your default cache policy here
    routes:
      - match:
          pathPrefix: assets/ # It's important not to start the path with a slash
        immutable: true
        public: true
        maxAge: 31536000s # 1 year
        sharedMaxAge: 31536000s # 1 year
        staleWhileRevalidate: 31536000s # 1 year
        staleIfError: 31536000s # 1 year
```

The result here is that all resources served by `spartan` with the path prefix
`assets/` (relative to the configured path root that is always normalized to
have a trailing slash, e.g. `/my-spa/`) will get a `cache-control` header like
this:

`Cache-Control: max-age=31536000, s-maxage=31536000, stale-if-error=31536000, stale-while-revalidate=31536000, immutable, public`

Routes are matched in the order in which they're configured. The first strategy
that results in the longest match is selected and its cache policy is applied.
For example, if a route matches `assets/index-` it will be preferred over a
route that matches just `assets/` because the former match is longer.

### Content security policy

Here's how the content security policy can be configured. The block must be
under `server.security` in the YAML file.

```yaml
contentSecurityPolicy:
  reportOnly: boolean 

  # Fetch directives; see https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#fetch_directives
  childSrc: fetch-directive
  connectSrc: fetch-directive
  defaultSrc: fetch-directive
  fencedFrameSrc: fetch-directive
  fontSrc: fetch-directive
  frameSrc: fetch-directive
  imgSrc: fetch-directive
  manifestSrc: fetch-directive
  mediaSrc: fetch-directive
  objectSrc: fetch-directive
  scriptSrc: fetch-directive
  scriptSrcElem: fetch-directive
  scriptSrcAttr: fetch-directive
  styleSrc: fetch-directive
  styleSrcElem: fetch-directive
  styleSrcAttr: fetch-directive
  workerSrc: fetch-directive

  # Document directives; see https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#document_directives
  baseUri: none-or-source-expression-list
  sandbox: all-or-list-of-allowed

  # Navigation directives; see https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#navigation_directives
  formAction: none-or-source-expression-list
  frameAncestors: none-or-source-expression-list

  # Reporting directives
  reportTo: string  # Defines the name of a reporting endpoint to send reports to.
```

The parameterless directives like `'none'`, `'self'`, `'unsafe-eval'` and so on
can all be configured using those exact strings, for example with `defaultSrc: "'none'"`
in YAML, but they can also be used without the extra single quotes for convenience:
`defaultSrc: none`.

#### Parametrized directives

The following directives with parameters can be used to further customize the
content security policy.

##### Nonce

See also [nonce-<nonce_value>](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#nonce-nonce_value).
Nonce directives are used to indicate a random value that establishes trust when
used with the `nonce` attribute of linked scripts and styles. `spartan`
generates a new nonce value every time a response is generated with a content
security policy that is configured with a nonce in a relevant policy directive.
It is configured as follows, for example for scripts and styles:

```yaml
scriptSrc:
  - nonce: PLACEHOLDER_VALUE_TO_REPLACE_WITH_NONCE
styleSrc:
  - nonce:
      placeholder: PLACEHOLDER_VALUE_TO_REPLACE_WITH_NONCE
```

The first form is shorter, the second form is more explicit. Both work and which
one you use is up to you. For the above example, `spartan` will generate _only
one_ nonce value for the placeholder `PLACEHOLDER_VALUE_TO_REPLACE_WITH_NONCE`
even when it is used in multiple directives. It will then replace the given
placeholder value in the responses with the generated nonce.

That is, if you have for example an HTML asset that references a script file
with a nonce as follows ...

```html
<html>
  <!-- ... -->
  <script type="module" crossorigin src="/assets/index-CZftpWSK.js" nonce="PLACEHOLDER_VALUE_TO_REPLACE_WITH_NONCE"></script>
  <!-- ... -->
</html>
```

... `spartan` will inject the generated nonce value into the `nonce` attribute
of the `<script>` element, using a different random value for every response.

##### Hash

See also [<hash_algorithm>-<hash_value>](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#hash_algorithm-hash_value).
Hash directives are used to indicate that a resource can be trusted if its hash
matches the hash configured in the content security policy. It is configured as
follows, for example for `<style>` elements:

```yaml
styleSrcElem:
  - hash: # The empty <style> element is trusted
      alg: sha384
      hash: OLBgp1GsljhM2TJ+sbHjaiH9txEUvgdDTAzHv2P24donTt6/529l+9Ua0vFImLlb
```

Allowed values for the `alg` (algorithm) properties are only `sha256`, `sha384`,
and `sha512`. The value for the `hash` property must be the base64 encoded value
of the hash for the trusted content/resource. For example, the `hash` configured
above represents the `sha384` hash value of empty content, thus in this case
explicitly trusting a style element that is empty: `<style></style>`.
You can verify this [here](https://gchq.github.io/CyberChef/#recipe=SHA2('384',64,160)From_Hex('Auto')To_Base64('A-Za-z0-9%2B/%3D')&oeol=FF).

##### Host

See also [<host_source>](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#host-source).
Host directives are used to indicate that resources served from a specific host,
with optional scheme, port and path, can be trusted. It is configured as follows,
for example for font resources:

```yaml
fontSrc:
  - host: fonts.gstatic.com
  - host:
      scheme: https
      host: my-fonts.com
      port: 443
      path: /trusted-fonts/
```

Here, fonts loaded from host `fonts.gstatic.com` (any protocol and any path) are
trusted, as are fonts loaded from host `my-fonts.com` when served over `https`
usings the default port of 443 and from a path under `/trusted-fonts/` on the
server.

##### Scheme

See also [<scheme_source>](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy#scheme-source).
Scheme directives are used to indicate that resources served from a protocol
(scheme) can be trusted. It is configured as follows, for example for connect
resources such as used when a script creates WebSockets:

```yaml
connectSrc:
  - scheme: wss
```

Here, scripts are allowed to open WebSocket connections as long as they are served
through a secured (web sockets over TLS) channel.

#### Default content security policy

If a content security policy is not configured, `spartan` defaults to using the
following configuration, expressed in YAML.

```yaml
server:
  security:
    contentSecurityPolicy:
      reportOnly: false
      defaultSrc: self
      objectSrc: none
      baseUri: none
      sandbox: all # This implies that all restrictions for sandboxing apply.
      formAction: self
      frameAncestors: self
```

This produces the following content security policy header in responses:

`Content-Security-Policy: default-src 'self'; object-src 'none'; base-uri 'none'; sandbox; form-action 'self'; frame-ancestors 'self'`

### Permissions policy

Here's how the permissions policy can be configured. The block must be under
`server.security` in the YAML file.

```yaml
permissionsPolicy:
  accelerometer: permissions
  ambientLightSensor: permissions
  ariaNotify: permissions
  attributionReporting: permissions
  autoplay: permissions
  bluetooth: permissions
  browsingTopics: permissions
  camera: permissions
  capturedSurfaceControl: permissions
  computePressure: permissions
  crossOriginIsolated: permissions
  deferredRetch: permissions
  deferredRetchMinimal: permissions
  displayCapture: permissions
  encryptedMedia: permissions
  fullscreen: permissions
  gamepad: permissions
  geolocation: permissions
  gyroscope: permissions
  hid: permissions
  identityCredentialGet: permissions
  idleDetection: permissions
  languageDetector: permissions
  localFonts: permissions
  magnetometer: permissions
  microphone: permissions
  midi: permissions
  onDeviceSpeechRecognition: permissions
  otpCredentials: permissions
  payment: permissions
  pictureInPicture: permissions
  publickeyCredentialsCreate: permissions
  publickeyCredentialsGet: permissions
  screenWakeLock: permissions
  serial: permissions
  speakerSelection: permissions
  storageAccess: permissions
  translator: permissions
  summarizer: permissions
  usb: permissions
  webShare: permissions
  windowManagement: permissions
  xrSpatialTracking: permissions
```

Where `permissions` can be one of `"*"`, `all`, or `wildcard` to grant the
permission indiscriminately. If you set it to `none` or `()`, the permission
is not granted at all.

You can grant a permission to more specific scopes by configuring a list. For
example, to grant the `geolocation` permission for several scopes:

```yaml
permissionsPolicy:
  geolocation:
    - self
    - src
    - https://my.origin.host.com
    - https://*.other.host.com
```

If you want to grant only a single scope for any permission, you can also
configure it directly without defining it as a list. For example, to allow `self`
for the `camera`:

```yaml
permissionsPolicy:
  camera: self
```

If the permissions policy is not configured, `spartan` defaults to using a
policy where each feature is not allowed, like shown
following configuration, expressed in YAML.

```yaml
server:
  security:
    permissionsPolicy:
      accelerometer: none
      ambientLightSensor: none
      ariaNotify: none
      # ...
```

This effectively implies that a permissions policy header like the following is
added to responses:

`Permissions-Policy: accelerometer=(), ambient-light-sensor=(), aria-notify=(), ...`

### Strict transport security policy

Here's how the strict transport security policy can be configured. The block
must be under `server.security` in the YAML file.

```yaml
# See the 'Duration values' section below to learn how to configure duration values.
strictTransportSecurityPolicy:
  disabled: boolean          # Set to true to disable the policy entirely
  includeSubDomains: boolean # Set to true to apply the policy to sub-domains too.
  maxAge: duration           # Defines how long browsers should remember the policy.
```

For example, the following strict transport security policy ...

```yaml
strictTransportSecurityPolicy:
  includeSubDomains: true
  maxAge: 8760h
```
... will produce a `strict-transport-security` header like
`Strict-Transport-Security: max-age=31536000; includeSubDomains`. This example
also represents the default policy unless it is overwritten with a different
policy or the `disabled` property is set to `true`.

### Duration values

The _duration_ values referenced above indicate a duration for the corresponding
directives and policies. These values are specified as a sequence of multiple
duration segments, each of which is a _number_ followed by a _unit_. Valid units
are `h` (hours), `m` (minutes), `s` (seconds).

For example, to set a _duration_ property to one year, set its value to `8760h`,
representing 365 * 24 hours. To set a _duration_ to 1 hour 2 minutes and 3
seconds, use the value `1h2m3s`.

## Attribution

The `spartan` logo was generated with Microsoft Copilot.
