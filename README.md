<img src="extension/images/page-ai-line.png" alt="logo" width="40" height="auto">

# ToRead

[![Website](https://img.shields.io/website?down_message=offline&up_color=green&up_message=online&url=https%3A%2F%2Ftr.lg.gl)](https://tr.lg.gl)

A simple Read-It-Later and link collection tool, AI-powered for text and images, multi-platform, open-source. A browser extension available for one-click bookmarking.

Demo deployed at <https://tr.lg.gl>

API document: <https://github.com/ligen131/ToRead/blob/main/backend/docs/api.md>

API Endpoint: <https://to-read.lg.gl/api/v1>

**Integrate ToRead into your application with just one line of code**:

```shell
$ curl https://g.lg.gl --data '{"url": "https://lg.gl"}'
```

The url field can be any publicly accessible link. Please note that some content requiring CAPTCHA verification or with restricted permissions may not be accessible. This endpoint may be slow to response.

This project was my undergraduate thesis at Huazhong University of Science and Technology (HUST). The current implementation is relatively simple and still has several bugs. The backend currently has no rate limiting or security checks, so please use it within a limited scope.

If you like this project, please give it a Star! ⭐

## Usage

### Browser Extension

Download the repository or just the extension folder separately.

Open your browser's extension settings, enable developer mode, and click "Load unpacked extension". Select the extension folder you downloaded.

Unhide the newly loaded ToRead extension from the top-right corner of your browser. Open any webpage, click the ToRead icon in the top-right corner to automatically add the link to your bookmarks. Note that you need to log in for first-time use.

## Deploy

### Backend

```shell
$ cp config-default.yml config.yml
```

Then mdify the configuration file `config.yml` to your needs.

Using golang 1.18+

```shell
$ go mod tidy
$ go run main.go
```

The backend will be running at `http://0.0.0.0:3435` by default.

### Frontend

```shell
$ npm install -g yarn
$ yarn
$ npm run start
```

The frontend will be running at `http://0.0.0.0:3000` by default.

## LICENSE

GNU General Public License v3.0
