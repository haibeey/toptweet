
var conntentString = `
<div class="jumbotron result-content fade-in">
    <div class="row">
        <h3 id="username" class="col-sm-9">user-name</h3>
        <a class="col-sm-3" href="link-to-tweet"><i class="fa fa-link">Link to Tweet</i></a>
    </div>
    <p id="the-twitter">con-tent</p>
    <div class="lrc">
        <i class="fa fa-comment lrc-items">cmnt</i>
        <i class="fa fa-retweet lrc-items">rtw</i>
        <i class="fa fa-heart lrc-items">hrt</i>
    </div>
</div>
`
var progressBar = `
<div id="p-boss">
<div class="progress" id="progress">
    <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" aria-valuenow="75" aria-valuemin="0" aria-valuemax="100" style="width: 100%"></div>
</div>
</div>`

var alertStr = `<div class="alert alert-danger" role="alert">
no-input
</div>`

function buildUrl(handle){
    return window.location.href+"search/?q=".concat(handle)
}

function processJson(responseText,handle){
    var result = JSON.parse(responseText)
    let divResult = document.getElementById("result")
    document.getElementById("p-boss").style.visibility = "hidden"
    for (let tweet of result.tweets){
        let divInner = document.createElement("div")
        divInner.innerHTML = conntentString.replace("user-name",handle).
                                            replace("con-tent",tweet.Text).
                                            replace("cmnt",tweet.ReplyCount).
                                            replace("rtw",tweet.RetweetCount).
                                            replace("hrt",tweet.FavoriteCount).
                                            replace("link-to-tweet",buildTweetURL(tweet.User,tweet.ID))
        divResult.appendChild(divInner)
    }
}
function fetchTopTweets() {
    document.getElementById("result").innerHTML = ""
    let pdiv = document.createElement("div")
    pdiv.innerHTML = progressBar
    document.getElementById("result").appendChild(pdiv)
    let inputElem = document.getElementById("search-input")

    if (inputElem.value){
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function() {
            if (this.readyState == 4){
                if (this.status == 200) {
                            processJson(this.responseText,inputElem.value)
                }else if(this.status == 400){
                    document.getElementById("p-boss").style.visibility = "hidden"
                    let divResult = document.getElementById("result")
                    let divInner = document.createElement("div")
                    divInner.innerHTML = alertStr.replace('no-input',JSON.parse(this.responseText).error)
                    divResult.appendChild(divInner)
                }
            }
        }
        xhttp.open("GET", buildUrl(inputElem.value), true);
        xhttp.send();
    }
}

function buildTweetURL(handle,id){
    return "https://twitter.com/".concat(handle,"/status/",id)
}