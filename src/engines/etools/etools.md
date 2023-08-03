# Etools

The first page request is a POST request that looks like:
https://www.etools.ch/searchSubmit.do
BODY: query=something&country=web&language=all&token=5d8d98d9a968388eeb4191afa00ca469 
Also works without token.

The requests for subsequent pages are GET requests that look like:
https://www.etools.ch/search.do?page=4
With a session cookie you got from some previous request: 
JSESSIONID=147933E3060CF19256C3581D55E7A72A

You can submit a GET request like: 
https://www.etools.ch/search.do?page=4&query=cool+cars
But you need the JSESSIONID cookie for it to work

It seems that, if performed too fast, the server can accidentaly return the same response for different pages. Thus the page requests are performed synchronously here. 


`?dataSourceResults=20` loads more requests

Possible settings to apply: `https://www.etools.ch/searchSettings.do`
Interesting are especially: `Results per search engine` and `Results per page`

Captcha Example:
![Alt text](captcha.png)