# Swisscows

Clicking search makes some HEAD, OPTIONS, GET requests:  
HEAD https://swisscows.com/_next/data/cW_SbMyHn51vQiG0qo8e9/en/web.json?query=cars on sale  
OPTIONS https://api.swisscows.com/web/search?query=cars+on+sale&offset=0&itemsCount=10&region=de-CH&freshness=All  
GET https://api.swisscows.com/web/search?query=cars+on+sale&offset=0&itemsCount=10&region=de-CH&freshness=All

We can use:
https://swisscows.com/en/web?query=\<query>&offset=\<(page-1) * 10>

Or we can directly call the API:  
https://api.swisscows.com/web/search?query=some+wrequest&offset=0&itemsCount=10&region=de-CH&freshness=All
Response:
![Alt text](image.png)

To use the API you have to pass these two headers:
![Alt text](image-1.png)

Example values (de-CH):
Pk_YSEsvfqugOyHxWFndYLtGzAoRLKM9
9hDmgLl6AS7wWH6PJCAPp4lGm1AGzw195HsPJS75qiU

The signature is dependant on the url parameters. So you can't use the same signature values for different request parameters. So, if you change the query, you have to get new signature values. 

## Reversing
function m(e, t) in line 11056 in app.js gone through a [beautifier](https://beautifier.io/) seems to assign these values.

Look at function R - GetAll in line 6856 in app.js

Sorted URL params are: ['freshness', 'itemsCount', 'offset', 'query', 'region']


URL: https://api.swisscows.com/web/search?query=something&offset=10&itemsCount=10&freshness=All
NONCE: 1YLVZK~znUCrDUxx~wf8gBqrObHG3otV
GIVES SIG: huCxYbxG1bMM6hYvAPUOXnUNZrS-MztSH2iTnJNnNBQ

params_dict = {
    "query":"something",
    "offset":"10",
    "itemsCount":"10",
    "freshness":"All"
}

sorted: 
/web/search?freshness=All&itemsCount=10&offset=10&query=something

I GENERATE: L3dlYi9zZWFyY2g_ZnJlc2huZXNzPUFsbCZpdGVtc0NvdW50PTEwJm9mZnNldD0xMCZxdWVyeT1zb21ldGhpbmcxbHlpbXh-TUFocEVxaEtLfkpTOFRvREViT3V0M0JHaQ

**THE ROT13 EMULATION WORKS.** The issue is probably in what `t.b64` does.

I need to get from:
/web/search?freshness=All&itemsCount=10&offset=10&query=something1lyimx~MAhpEqhKK~JS8ToDEbOut3BGi
to
huCxYbxG1bMM6hYvAPUOXnUNZrS-MztSH2iTnJNnNBQ

`t.b64` is this.b64 on line 5361 in app.jsnice.js

SHA256 of /web/search?freshness=All&itemsCount=10&offset=10&query=something1lyimx~MAhpEqhKK~JS8ToDEbOut3BGi is:
86e0b161bc46d5b30cea162f00f50e5e750d66b4be333b521f68939c93673414

[SHA256 JS repo](https://gist.github.com/napkindrawing/758673)

site/tcopy.js and site/tcopy.nice.js work perfectly.

Converting the complicated JS to Golang is hard, so we just run a JS parser in `dontaskjustenjoy.go`.

### Parsing
The retrieved json results have HTML tags in them. The current way of escaping them could be improved.