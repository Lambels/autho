# Autho
Autho is a customisable and clean implementation of all OAuth and OAuth2.0 web providers with full implementation for popular providers such as [twitter](https://github.com/Lambels/autho/tree/main/twitter), [facebook](https://github.com/Lambels/autho/tree/main/facebook), [github](https://github.com/Lambels/autho/tree/main/github) and [google](https://github.com/Lambels/autho/tree/main/google). Autho helps you focus on the logic side of your application whilst providing an idiomatic api to handle the boilerplate.

Autho was inspired by: [gologin](https://github.com/dghubble/gologin) and [goth](https://github.com/markbates/goth)

# Quickstart 🚀
## Install autho:
```
go get github.com/Lambels/autho
```

## Github Provider:
```go
import (
    "github.com/Lambels/autho"
    gh "github.com/Lambels/autho/github"
	"github.com/google/go-github/v32/github"
    "golang.org/x/oauth2"
	oauthGh "golang.org/x/oauth2/github"
)

func main() {
    ckCfg := autho.NewDebugCookieConfig("my-cookie")
    ghCfg := &oauth2.Config{
        ClientID: "client-ID",
        ClientSecret: "client-secret",
        Endpoint: oauthGh.Endpoint,
    }
    app := autho.NewApp(
        autho.NewProviderApp(
            "github/login",
            "github/callback",
            gh.NewCallbackHandler(ghCfg, ckCfg, nil, terminalHandler),
            gh.NewLoginHandler(ghCfg, ckCfg),
        ),
    )
    mux := http.NewServeMux()
    app.Register(mux)
    srv := &http.Server{
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

// this is what gets called at the end of the OAuth process if everything goes right.
// use the request context to pull the *github.User .
func terminalHandler(w http.ResponseWriter, r *http.Request) {
    // at this point we for sure have a user object under the ctx
    // but still check for ok to avoid any panics.
    user, ok := autho.UserFromContext(r.Context()).(*github.User)
    if !ok {
        return
    }

    fmt.Println(user)
}
```

# Docs
[GoDoc](https://pkg.go.dev/github.com/Lambels/autho)

# Overview (Callback and Login Phases)
Simillarly to [gologin](https://github.com/dghubble/gologin), `autho` chains `http.Handler`s to build handlers for the two main stages of OAuth flows: "The Login Phase" and "The Callback Phase", in `autho` the login phase has its own handlers in the `autho/oauth1` and `autho/oauth2` packages which are at the core of autho. `autho/oauth1.NewLoginHandler()` and `autho/oauth2.NewLoginHandler()` are the two core login handlers, then packages such as `autho/facebook.NewLoginHandler()` or `autho/twitter.NewLoginHandler()` provide abstractions of those handlers respectivly for top level use. Simillarly the callback phase has its own handler which is constructed by wrapping the providers specific user handler with the protocols token handler.

Facebook CallbackHandler example: [Oauth2 Token Handler](https://github.com/Lambels/autho/blob/main/oauth2/oauth2.go#L46) + [Facebook User Handler](https://github.com/Lambels/autho/blob/main/facebook/oauth2.go#L54) = [Facebook Callback Handler](https://github.com/Lambels/autho/blob/main/facebook/oauth2.go#L18)

# LoginHandler OAuth1.0 And OAuth2.0
The login handler is in both OAuth1.0 and OAuth2.0 responsible for redirecting the user to the provider for obtaining the users grant. However there are differences:

OAuth1.0:
- The [OAuth1.0 Login Handler](https://github.com/Lambels/autho/blob/main/oauth1/oauth1.go#L17) is responsible for, if the provider requires it, setting a cookie with the request secret for later reading in the Token Handler. This behaviour is flagged by passing `nil` to the cookie config parameter, if its nil the provider doesent require it, if it isnt the provider requires and the cookie must be set accordingly.

OAuth2.0:
- The [OAuth2.0 Login Handler](https://github.com/Lambels/autho/blob/main/oauth2/oauth2.go#L18) is responsible for setting the state value in a short lived cookie to be validated in the token handler step.

# CallbackHandler
The callback handler is specific to each provider and can be built either by steps or by using the providers helper method.

`autho/twitter.NewCallbackHandler()` or `autho/github.NewCallbackHandler()` are examples of building callback handlers using the provided helper methods in each provider package. In essence what these helper functions are doing is using the default token and user handlers for each provider to build the callback handler.

Basically a callback handler is just an `http.Handler` but obtained through chaining multiple `http.Handler`s.

The flow looks something like this:

OAuth2.0:

Provider (redirects to your server) -> TokenHandler (validates state + exchanges auth code for tokens) -> UserHandler (uses obtained tokens to exchange for user. unique to each provider) -> TerminalHandler (has access to tokens and user resource. end logic)

OAuth1.0:

Provider (redirects to your server) -> TokenHandler (exchanges request token + request secret + verifier for tokens) -> UserHandler (uses obtained tokens to exchange for user. unique to each provider) -> TerminalHandler (has access to tokens and user resource. end logic)

The main difference between OAuth2.0 and OAuth1.0 is in the token hanlder.

## TokenHandler OAuth1.0 And OAuth2.0
The token handler is in both OAuth1.0 and OAuth2.0 responsible for, as the name says, obtaining the tokens used in exchange for the users resource. However there are some differences:

OAuth1.0:
- The [OAuth1.0 Token Handler](https://github.com/Lambels/autho/blob/main/oauth1/oauth1.go#L60) is responsible for grabbing the request secret from the short lived cookie if the provider requires so. Some providers require that the request secret is persisted throughout the exchange, some dont. This behaviour is flagged by passing `nil` to the cookie config parameter, if its nil the provider doesent require it, if it isnt the provider requires and the cookie must be read accordingly.

OAuth2.0:
- The [OAuth2.0 Token Handler](https://github.com/Lambels/autho/blob/main/oauth2/oauth2.go#L46) is responsible for validating the state from the short lived cookie and compare it with the state from the request.

Finally both token handlers add the tokens to the request context to be used down the line by the user handler using the `autho/oauth1.ContextWithToken()` or `autho/oauth2.ContextWithToken()` respectively.

## UserHandler
The user handler is specific to each provider and comes in chain right after the token handler, it is responsible for exchanging the tokens obtained from the request context using the `autho/oauth1.TokenFromContext()` or `autho/oauth2.TokenFromContext()` for the user resource.

Finally after obtaining the user resource the user handler adds the user object to the request context using the `autho.UserWithContext()` method, it then calls the terminal handler which is last in chain.

## TerminalHandler
The terminal handler is the "end logic" and must be implemented by you. To access the tokens if your provider uses OAuth1.0 use `autho/oauth1.TokenFromContext()` else use `autho/oauth2.TokenFromContext()`. To access the user resource use `autho.UserFromContext()` which returns `interface{}` so it is up to you to parse the `interface{}` to your own type. To know what type the user is of check your providers user handler docs which specifies the user type.

## ErrorHandler
Obviously throughout the whole OAuth1.0 or OAuth2.0 flow errors can occur, the error handler gets called by handlers when an error occurs.

This is githubs callback handler signature:

`autho/github.NewCallbackHandler(cfg *oauth2.Config, ckCfg *autho.CookieConfig, errHandler, terminalHandler http.Handler)`

The `errHandler` parameter is a `http.Handler`, if `nil` is passed the `autho.DefaultFailureHandler` is used, else you can implement your own. The error inside the handler is obtainable by the request context using the `autho.ErrorFromContext()` and its up to you how you handle it.

# Customising The Handlers
There are essentially 5 `http.Handler`s in the whole exchange, but you can chain as many as you want by chaining n `http.Handler`s.

The 5 handlers used by `autho` defaultly are:
- LoginHandler
- TokenHandler
- UserHandler
- TerminalHandler
- ErrorHandler

The flow is:
LoginHandler -> Provider
Provider -> TokenHandler -> UserHandler -> TerminalHandler (any errors -> ErrorHandler)

Obviously you can provide your own LoginHandler or CallbackHandler as long as it is an `http.Handler` to `autho.NewProviderApp()`

When you call `autho/someProvider.NewCallbackHandler()` the helper function chains:
```
TokenHandler(
    ErrorHandler,
    UserHandler(
        ErrorHandler,
        TerminalHandler,
    )
)
```

# Code Examples

## Github OAuth2.0
This is a full github provider implemented, the `ghCfg` is a `*golang.org/x/oauth2.Config` and the `ckCfg` is a `*autho.CookieConfig`. The terminal handler is the end logic which must be implemented by the user. The user object is obtainable from the request context with the method `autho.UserFromContext()`, to know what type to parse the user object to, check the providers user handler docs, in our case its the [Github User Handler](https://github.com/Lambels/autho/blob/main/github/oauth2.go#L54).

```go
app := autho.NewApp(
    autho.NewProviderApp(
        "github/login",
        "github/callback",
        gh.NewCallbackHandler(ghCfg, ckCfg, nil, terminalHandler),
        gh.NewLoginHandler(ghCfg, ckCfg),
    ),
)

func terminalHandler(w http.ResponseWriter, r *http.Request) {
    // at this point we for sure have a user object under the ctx
    // but still check for ok to avoid any panics.
    user, ok := autho.UserFromContext(r.Context()).(*github.User)
    if !ok {
        return
    }

    fmt.Println(user)
}
```

## Twitter OAuth1.0
Simillarly to the Github implementation, we will use the `autho.NewProviderApp()` method to create our twitter provider. Only difference here is the `twCfg` and the source of the handlers. The `twCfg` is a `*github.com/dghubble/oauth1.Config`. The Login and Callback hanlder come from the `autho/twitter` package. A complete abstraction is made between OAuth1.0 and OAuth2.0 .

```go
app := autho.NewApp(
    autho.NewProviderApp(
        "twitter/login",
        "twitter/callback",
        tw.NewCallbackHandler(twCfg, ckCfg, nil, terminalHandler),
        tw.NewLoginHandler(twCfg, ckCfg),
    ),
)

func terminalHandler(w http.ResponseWriter, r *http.Request) {
    // at this point we for sure have a user object under the ctx
    // but still check for ok to avoid any panics.
    user, ok := autho.UserFromContext(r.Context()).(*twitter.User)
    if !ok {
        return
    }

    fmt.Println(user)
}
```