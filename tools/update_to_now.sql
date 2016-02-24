-- update audience and article to now
update audience set crawl_time = crawl_time + (select now() - crawl_time from audience order by crawl_time desc limit 1);
update article set crawl_time = crawl_time + (select now() - crawl_time from article order by crawl_time desc limit 1);
