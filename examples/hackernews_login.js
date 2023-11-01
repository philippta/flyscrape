import { submitForm } from "flyscrape"

export const config = {
    url: "https://news.ycombinator.com",
}

export function login() {
    const formData = {
        "acct": "my-username",
        "pw": "my-password",
    }

    submitForm("https://news.ycombinator.com/login", formData)
}

export default function ({ doc }) {
    return {
        karma: doc.find("#karma").text()
    }
}
