# Presearch

It's open source, but there doesn't seem to be any website code: https://github.com/PresearchOfficial

GET request: https://presearch.com/search?q=something&page=3  
Gets populated with API call: GET https://presearch.com/results?id=5b747ca66cc051a82a6c5bbb784a7fa5f802

There are cookies:
+ settings cookies:
  + ai_results_disable:1
  + use_safe_search:true
+ session cookies:
  + presearch_session: eyJpdiI6InBtNVgzZE5YZnUvcXRldGNrZytzTWc9PSIsInZh[...]
  + XSRF-TOKEN: eyJpdiI6InN5MlM1Z3ovdkJuQzNBcW5MM0x6RkE9PSIsInZhbHVlI[...]
+ weird cookies:
  + b: 0
  + AWSALB: N5A3Uv4njhnPnihhwOzEBPWXwUZCx/KyphsluMdnYHL[...]
  + AWSALBCORS: N5A3Uv4njhnPnihhwOzEBPWXwUZCx/KyphsluMdnY[...]

The id to pass to results is the JS variable "window.searchId" that gets set on the initial GET request, it is generated server-side