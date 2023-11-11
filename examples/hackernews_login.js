import http from "flyscrape/http"

export const config = {
    url: "https://news.ycombinator.com",
}

export function setup() {
    http.postForm("https://news.ycombinator.com/login", {
        "acct": "my-username",
        "pw": "my-password",
    })
}

export default function ({ doc }) {
    return {
        karma: doc.find("#karma").text()
    }
}
