# GMC VIEW

gmc view module, make it very easy to render your template.

## CONFIGURATION

about template configuration is section [template] in app.toml.

```toml
[template]
dir="views"
ext=".html"
delimiterleft="{{"
delimiterright="}}"
```

## FUNCTIONS

### text/template inside functions

`and`, `call`, `html`, `index`, `slice`, `js`, `len`, `not`, `or`, `print`, `printf`, `println`, `urlquery`

### text/template inside comparison functions

`eq`, `ne`, `lt`, `le`, `gt`, `ge`

### gmc defined functions

1. `tr` be used in i18n, first argument is always `.Lang`, secondary is key name in your i18n locale config file.
    third is optional tips text for yourself , it will be not output. returns template.HTML typed string.
    
    Example in the template:
    
    `{{tr .Lang "key001" "tips"}}`

1. `trs` as same as `tr`, but returns string type.

1. `string` only one argument, type is `interface{}`, convert it to string type. 

1. `tohtml` only one argument, convert it to template.HTML type.

### gmc defined comparison functions


### gmc defined variables

1. `.Lang` is the i18n standard FLAG the result of parsed client browser's Accept-Language HTTP header.  
    It is be used in the i18n `tr` function.
1. `.G` access GET data from URL query string. `.G` is a `map[string][string]`, 
    `.G.key` key is query name in query string.
1. `.P` access POST form data. `.P` is a `map[string][string]`, 
    `.P.key` key is form field name in POST form.
1. `.S` access session data from the current session. `.S` is a `map[string][string]`, 
    `.S.key` key is query name in query string.  
    Notice that this only worked after the SessionStart() in a controller is be called.
1. `.C` access COOKIE data. `.C` is a `map[string][string]`, 
    `.C.key` key is cookie field name in cookies raw string.
1. `.U` current URL information `u:=map[string]string`.
    
    ```golang
    u["HOST"] = u0.Host
    u["HOSTNAME"]=u0.Hostname()
    u["PORT"]=u0.Port()
    u["PATH"] = u0.Path
    u["FRAGMENT"] = u0.Fragment
    u["OPAQUE"] = u0.Opaque
    u["RAW_PATH"] = u0.RawPath
    u["RAW_QUERY"] = u0.RawQuery
    u["SCHEME"] = u0.Scheme
    u["USER"] = u0.User.Username()
    u["PASSWORD"],_ = u0.User.Password()
    u["URI"]=u0.RequestURI()
    u["URL"]=u0.String()
    ```
   TIPS: `[scheme:][//[userinfo@]host][/]path[?query][#fragment]`  



