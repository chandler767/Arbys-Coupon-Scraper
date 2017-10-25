# Arby's Coupon Scraper
This is a concurrent web page scraper built in Go for finding coupons for Arby's fast food restaurants. This takes advantage of the API Arby's uses to display coupons for when customers need to print them from their email. This code could easily be reused for another project and demonstrates some of the capabilities of the Go language. 

#### How it works:
1. Sign up for [Arby's email deals](https://arbys.com/get-deals).
2. Wait for a coupon to be emailed to you.
3. Click where it says "Click here for a printable version".
4. You should get a URL like this: `http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=27917361605`
5. Remove any query values after the MailingID in the URL.
  - Notice that if you increase or decrease the MailingID value you can get a different coupon.
  - The scraper works by starting at the MailingID you specify and increases the MailingID until it hits the total requests value.
  - Each valid page is saved and the "Offer valid only at:" text and expiration dates are removed.
  
### How to use it:
1. Set the variables to your liking.
2. `go run main.go`

  ```
  flag.IntVar(&reqs, "reqs", 500, "Total requests")                                                                     // How many MailingID vaules you want to check. This scraper only scrapes values higher than this one. If you want old or more coupons use a lower number for the MailingID and increase the total requests.
  flag.IntVar(&max, "concurrent", 45, "Maximum concurrent requests")                                                    // Prevents "Too many open files" error and timeouts.
  flag.IntVar(&mailingid, "mailingid", 27917361605, "Mailing ID to start on")                                           // Get a recent MailingID from an Arby's promo email by clicking to view print version and getting the query value for 'MailingID' from the url. URL Example from 10/25/2017: http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=27917361605
  flag.StringVar(&scrapeurl, "scrapeurl", "http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=", "URL to scape") // This was found by viewing a printable version of an Arby's promo email and removing the extra query values from the URL.
  flag.BoolVar(&formatcoupon, "formatcoupon", true, "Remove the 'Offer valid only at:' text and expiration dates.")     // Remove the "Offer valid only at:" text and expiration dates.
  ```

 3. You should get an output like this:
  ```
  ...
  2017/10/25 11:14:16 http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=27917361938
  2017/10/25 11:14:16 http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=27917361944
  2017/10/25 11:14:16 http://arbys.fbmta.com/members/ViewMailing.aspx?MailingID=27917361945
  Connections:	500
  Concurrent:	45
  Total time:	15.796205411s
  Average time:	31.59241ms
  Total results:	40
  ```
  4. HTML files are saved to the path where main.go is located. 
