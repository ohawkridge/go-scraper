# _GO_ Web Scraper
[LIVE DEMO](https://some.com/) 🔗

### What does it do?

- INPUT ← Scrapes a popular teaching jobs website and stores the jobs in MySQL.
- PROCESSING — Basic data processing like parsing dates.
- OUTPUT → Folders of static HTML files with different views into jobs. For example, one folder contains jobs organised by location.

### Architecture

- Web scraper and static file generation use GO and run on local machine.
- Database uses MySQL hosted on AWS RDS.
- Static files hosted on Netlify. These are all pre-generated so browsing the output is lightening fast ⚡️.
