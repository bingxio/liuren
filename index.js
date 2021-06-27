// index.js
// 
// @author: bingxio
// @create: 2021/06/26 13:40
// 
let earthBranch = ['子', '丑', '寅', '卯', '辰', '巳', '午', '未', '申', '酉', '戌', '亥']
let heavenStems = ['甲', '乙', '丙', '丁', '戊', '己', '庚', '辛', '壬', '癸']

let goTo = () => {
    let a = document.getElementById('a')
    let b = document.getElementById('b')

    let c = document.getElementById('t2')
    let d = document.getElementById('t3')

    let t1 = a[a.selectedIndex].innerText + b[b.selectedIndex].innerText
    let t2 = c[c.selectedIndex].innerText
    let t3 = d[d.selectedIndex].innerText

    console.log(t1[0])
    console.log(t1[1])

    let i = 0
    for (; i < heavenStems.length; i ++) {
        if (heavenStems[i] == t1[0]) {
            i ++
            break
        }
    }

    let j = 0
    for (; j < earthBranch.length; j ++) {
        if (earthBranch[j] == t1[1]) {
            j ++
            break
        }
    }

    if (
        i % 2 == 0 && j % 2 != 0 || i % 2 != 0 && j % 2 == 0
    ) {
         window.alert('请选择匹配的正确干支')
         return
    }

    window.location.href = './paipan.html?x=' 
        + t1 + '&y=' + t2 + '&z=' + t3
}