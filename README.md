# Spartan

> A simple and secure web server for SPA or similar apps with static assets.

## Usage

```
spartan [flags]

Flags:
  -p, --port uint16                 The local port to listen on for incoming requests. (default 8080)
  -r, --server-path-root string     The absolute path on the server where the static content is exposed.
      --config string               The config file to use.
      --logtostderr                 log to standard error instead of files (default true)
  -d, --static-content-dir string   The path to the directory holding the static content to serve. (default "/content")
  -v, --v Level                     number for the log level verbosity
  -h, --help                        help for spartan
```

For example, the follogin command starts the `spartan` web server for static
content from the directory `/path/to/web` and exposing it on the sub-path `/web/`
through the HTTP interface:

```bash
spartan -d /path/to/web -r /web/
```

That is, when a browser requests for example _http://localhost:8080/web_, the
browser will automatically render the file _/path/to/web/index.html_ served by
`spartan`.

### Container usage

When using the official spartan container images, `spartan` is located at
`/srv/spartan` and it expects the `config.yaml` file in the same directory, i.e.,
it tries to load `/srv/config.yaml`. However, you can of course put the YAML
configuration file in any location of your liking and instruct `spartan` to load
it from there. For example, with _docker_:

```bash
docker container run --rm -it spartan --config=/path/to/my-config.yaml
```

Please note that of course you should **not** put the configuration file in the
same folder (or a subfolder) as the static content, since otherwise your
configuration file will be exposed by `spartan` itself.

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

  cache:
    defaultPolicy:
      # The cache policy to use, see below.

  security:
    contentTypeOptionsNoSniff: true
    referrerPolicy: strict-origin-when-cross-origin
    contentSecurityPolicy:
      # The content security policy to use, see below.
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
