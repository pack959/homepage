# Pack 959 [![Netlify Status](https://api.netlify.com/api/v1/badges/d115f79f-e3b4-419e-b682-4e9b343eec7a/deploy-status)](https://app.netlify.com/sites/pack959/deploys)

This is the homepage for Pack 959 (https://pack959.com/).

## Local Development

1. [Install Hugo](https://gohugo.io/overview/installing/)
1. Clone this repository

    ```bash
    git clone https://github.com/pack959/homepage.git
    cd homepage
    ```

1. Run Hugo test server

    ```bash
    hugo server
    ```

1. Generate calendar events

    ```bash
    make genevents
    ```

    You can configure the date range for the calendar events in the Makefile by
    setting `CALSTART` and `CALEND` variables.