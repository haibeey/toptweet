var conntentString = `
<div class="jumbotron result-content">
<h3 id="username">user-name</h3>
<p id="the-twitter">con-tent</p>
<div class="lrc">
    <i class="fa fa-comment lrc-items">cmnt</i>
    <i class="fa fa-retweet lrc-items">rtw</i>
    <i class="fa fa-heart lrc-items">hrt</i>
</div>
</div>
`
function buildUrl(handle){
    return window.location.href+"search/?q=".concat(handle)
}

function processJson(responseText,handle){
    var result = JSON.parse(responseText)
    let div = document.getElementById("result")
    for (let tweet of result.tweets){
        let divInner = document.createElement("div")
        divInner.innerHTML = conntentString.replace("user-name",handle).
                                            replace("con-tent",tweet.Text).
                                            replace("cmnt",tweet.ReplyCount).
                                            replace("rtw",tweet.RetweetCount).
                                            replace("hrt",tweet.FavoriteCount)
        div.appendChild(divInner)
    }
}
function fetchTopTweets() {

    let inputElem = document.getElementById("search-input")
    console.log(inputElem.value,"what",buildUrl(inputElem.value))
    if (inputElem.value){
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
                    processJson(this.responseText,inputElem.value)
            }
        }
        xhttp.open("GET", buildUrl(inputElem.value), true);
        xhttp.send();
    }
}
