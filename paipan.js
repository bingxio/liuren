// paipan.js
// 
// @author: bingxio
// @create: 2021/06/26 13:58
// 
window.onload = () => {
    let t1 = getQuery('x')
    let t2 = getQuery('y')
    let t3 = getQuery('z')

    let xhr = null

    if (window.XMLHttpRequest) {
        xhr = new XMLHttpRequest()
    } else {
        xhr = new ActiveXObject('Microsoft.XMLHTTP')
    }

    let url = 'http://127.0.0.1:8080/get?x=' + t1 + '&y=' + t2 + '&z=' + t3
    let body = document.getElementById('body')

    xhr.open('GET', url, true)
    xhr.onreadystatechange = () => {
        if (xhr.readyState == 4 && xhr.status == 200) {
            let data = JSON.parse(xhr.response).data
            let tb = document.createElement('table')

            tb.innerHTML = data.x
            body.appendChild(tb)

            let tb2 = document.createElement('table')
            tb2.innerHTML = data.y
            body.appendChild(tb2)
        }
    }
    xhr.send(null)
}

let getQuery = (name) => {
    let reg = new RegExp('(^|&)' + name + '=([^&]*)(&|$)')
    let r = decodeURI(window.location.search.substr(1)).match(reg)

    if (r != null) {
        return unescape(r[2])
    } else {
        return null
    }
}