# multilogin_scraping

# How to run the script
1. Open port 35000 on your Multilogin App by following this guide: https://docs.multilogin.com/l/en/article/el0fuhynnz-a-quick-guide-to-starting-browser-automation
2. Edit your config file with main options such as:
      - crawler.zillow_crawler.time_load_source: Time to read source page from browser.
      - crawler.zillow_crawler.crawl_next_time: Time to crawl next address
   1. Crawl Zillow Data
      - crawler.zillow_crawler.periodic_record_size: Total addresses / number of browsers => The browser will stop if arrive to this value 
      - crawler.zillow_crawler.periodic_browser: number of browsers
      - crawler.zillow_crawler.periodic_run: Next time to open new browsers
   2. Crawl price and tax histories
      - crawler.zillow_crawler.periodic_record_size_interval: Total addresses / number of browsers => The browser will stop if arrive to this value
      - crawler.zillow_crawler.periodic_browser_interval: number of browsers
      - crawler.zillow_crawler.periodic_interval: Next time to open new browsers
      - crawler.zillow_crawler.days_interval: How many days will be crawling data again.

4. Open your Windows's terminal and run command ./multilogin_scraping.exe  