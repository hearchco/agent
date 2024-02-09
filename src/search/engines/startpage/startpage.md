# Startpage

First search: POST request to https://www.startpage.com/sp/search
with body: abp=-1&additional=%5Bobject+Object%5D&cat=web&language=english&lui=english&query=some+query&sc=BSuId774jcrp20&sgt=1691175704T0afc510362af195aa4ac76bde15e32e85914a4901124669719eaac0e2c326f15&t=

Sending just cat,language,lui,query gets this:
![Alt text](image.png)

Resending the previous request gets this:
![Alt text](image-1.png)

Request to second page: POST request to https://www.startpage.com/sp/search
with body: language=english&lui=english&abp=-1&query=some+query&cat=web&page=2&sc=HLlIFdefdQOM20

Resending it worked fine.

Changing HLlIFdefdQOM20 to HLlIFdefdZOM20 and resending worked fine. Changing it to aaaaaaaaaaaaaa redirects to an error page, that sends the javascript message. The sc value is plainly set in the html (form#search > input[name="sc"]). When last page is hit:
![Alt text](image-2.png)

Doesnt use cookies.

+ Safe search is on: add qadf=heavy to POST body
+ Safe search is off: add qadf=none to POST body
- Not sure if it needs to be set with every request

Disabling javascript in browser settings gets the **Error 883** page. However, sending requests through GET: https://www.startpage.com/sp/search?q=<query> works even if javascript is disabled. The GET request works with no cookies / body. For the page, the `page` URL parameter is used. E.g. https://www.startpage.com/sp/search?q=i+dont+get+it&page=3


# Locale
The locale is set with the POST body `qloc` variable and looks something like this:
`JTdCJTIyY2MlMjIlM0ElMjJVUyUyMiUyQyUyMmxvY2F0aW9uJTIyJTNBJTIyVW5pdGVkJTIwU3RhdGVzJTIyJTJDJTIyc3RhdGVfY29kZSUyMiUzQSUyMjAwJTIyJTJDJTIydHlwZSUyMiUzQSUyMmN1c3RvbV9sb2NhdGlvbiUyMiU3RA%3D%3D`\
If we replace `%3D`s with `=`s we can base64 decode it into:
`%7B%22cc%22%3A%22US%22%2C%22location%22%3A%22United%20States%22%2C%22state_code%22%3A%2200%22%2C%22type%22%3A%22custom_location%22%7D`\
Which we can url decode into:
`{"cc":"US","location":"United States","state_code":"00","type":"custom_location"}`
Another decoded example is
`{"cc":"CN","location":"People’s Republic of China","state_code":"00","type":"custom_location"}`
It seems for states, the `state_code` is always `00` and the `type` is always `custom_location`. The `location` parameter may be irrelevant, and spoofing `cc` could be sufficent. However, the results dont seem to change when the region is changed, so its impossible to test.
