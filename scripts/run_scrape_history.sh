#!/bin/bash

# NHL
./oddsportal-scraper -m "full" -u https://www.oddsportal.com/hockey/usa/nhl-2018-2019/results/#/page/ -s "./data/NHL/20182019/" -f "./data/NHL/20182019" -o true
./oddsportal-scraper -m "full" -u https://www.oddsportal.com/hockey/usa/nhl-2019-2020/results/#/page/ -s "./data/NHL/20192020/" -f "./data/NHL/20192020" -o true
./oddsportal-scraper -m "full" -u https://www.oddsportal.com/hockey/usa/nhl-2020-2021/results/#/page/ -s "./data/NHL/20202021/" -f "./data/NHL/20202021" -o true
./oddsportal-scraper -m "full" -u https://www.oddsportal.com/hockey/usa/nhl-2021-2022/results/#/page/ -s "./data/NHL/20212022/" -f "./data/NHL/20212022" -o true
./oddsportal-scraper -m "full" -u https://www.oddsportal.com/hockey/usa/nhl-2022-2023/results/#/page/ -s "./data/NHL/20222023/" -f "./data/NHL/20222023" -o true
./oddsportal-scraper -m "full" -u https://www.oddsportal.com/hockey/usa/nhl-2023-2024/results/#/page/ -s "./data/NHL/20232024/" -f "./data/NHL/20232024" -o true
