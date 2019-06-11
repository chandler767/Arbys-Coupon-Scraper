# Arby's Coupon Scraper
This is a scraper built in Go that finds coupons for Arby's fast food restaurants. This takes advantage of the url formatting Arby's uses to display coupons for when customers need to print them from their email.

#### How it works:
1. Sign up for [Arby's email deals](https://arbys.com/get-deals).
2. Wait for a coupon to be emailed to you.
3. Click where it says "Click here for a printable version".
4. You should get a URL like this: `http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=27917361605`
5. Remove any query values after the MailingID in the URL.
  - Notice that if you increase or decrease the MailingID value you can get a different coupon.
  - The scraper works by starting at the MailingID you specify and increases the MailingID until it hits the total requests value.
  - Each valid page is saved and the "Offer valid only at:" text and expiration dates are removed be default.
  
### How to use it:
1. `go run main.go` or `go run main.go -id 27917361605`
  ```
  Usage:
  -concurrent int
    	Maximum concurrent requests (default 45)
  -format bool
    	Remove the 'Offer valid only at:' text and expiration dates. (default true)
  -id int
    	Mailing ID to start on (default 27917361605)
  -total int
    	Total requests (default 500)
  -url string
    	URL to scape (default "http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=")
  ```
 2. You should get an output like this:
  ```
  ...
  2019/06/11 12:12:23 http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=27917362940
  2019/06/11 12:12:23 http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=27917362958
  2019/06/11 12:12:24 http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=27917363016
  Connections:	500
  Concurrent:	45
  Total time:	2.796205411s
  Average time:	5.59241ms
  Total results:	30
  ```
  3. HTML files are saved to the current working directory. 
