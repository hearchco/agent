# Qwant

We access the api (https://api.qwant.com/v3/search/web) and set the necessary headers: <br>
```
q: <query>
count: 10
locale: en_GB
offset: 10
device: desktop
safesearch: 1
```

To parse the incoming JSON we use https://pkg.go.dev/encoding/json#Unmarshal ([help](https://www.sohamkamani.com/golang/json/)). Especially note:
> By default, object keys which don't have a corresponding struct field are ignored (see Decoder.DisallowUnknownFields for an alternative).

We pass data to the colly callbacks like this:
```
colCtx := colly.NewContext()
colCtx.Put("offset", strconv.Itoa(i*qResCount))
col.Request("GET", seURL, nil, colCtx, nil)
```
^ Instead of colly.Visit(seURL)

For the first result page `col.Visit(seURL + query + "&t=web&locale=" + qLocale + "&s=" + qSafeSearch)` could be used. This would emulate an actual user better. Its `.OnHTML` is implemented, but it seems to not play well with the API calls, having some results overlapp, this doesn't make any sense whatsoever. If this is used for first page, then `for i := 0; i < options.MaxPages; i++ {` needs start at 1 (i.e. `for i := 0; ....`). When it works and when it doesn't seems random - so it may be best to not touch it. Last query on which it didn't work: `./main --query="jako cudne stvari" --max-pages=2 -vv --visit`