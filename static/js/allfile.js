function updatepage(json_lst){
    var gallery = document.getElementById("gallery");
    var ul = gallery.querySelector("ul")
    const liElements = ul.querySelectorAll('li'); // 获取所有 <li> 子元素

    liElements.forEach(li => ul.removeChild(li)); 
    if (json_lst.length>0){
        json_lst.forEach(item => {
            const img = document.createElement('img');
            const span = document.createElement("span");
            const aremove = document.createElement("a");
            const aview = document.createElement("a");
            // 设置li
            const li = document.createElement("li");
            li.className = "item";
            li.id = item["FileId"]
            // 设置a span img
            img.className = "item-image";
            img.alt = "Placeholder Image";
            img.src  = "./file/" + item["FileId"];
            span.className = "item-name";
            span.textContent =item["Name"];
            aremove.className = "custom-button";
            aremove.textContent = "删除";
            aremove.addEventListener('click', (event) => removeimage(event,item["FileId"]));
            aview.className = "custom-button";
            aview.textContent="查看";
            aview.href = "./file/" + item["FileId"];
            li.append(img)
            li.append(span)
            li.append(aview)
            li.append(aremove)
            ul.append(li)
        })
    }else{
        const img = document.createElement('img');
        const span = document.createElement("span");
        const aremove = document.createElement("a");
        const aview = document.createElement("a");
        // 设置li
        const li = document.createElement("li");
        li.className = "item";
        // 设置a span img
        img.className = "item-image";
        img.alt = "Placeholder Image";
        img.src  = "./file/"
        span.className = "item-name";
        span.textContent ="无文件内容";
        aremove.className = "custom-button";
        aremove.textContent = "删除";
        aview.className = "custom-button";
        aview.textContent ="查看"
        li.append(img)
        li.append(span)
        li.append(aview)
        li.append(aremove)
        ul.append(li)
    }

}

async function removeimage(event,fileid){
    event.preventDefault(); // 阻止默认行为（如跳转）
    var formData = new FormData();
    formData.append("fileid",fileid)
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/remove", true);
    xhr.onload = function (){
        var json_data = JSON.parse(xhr.responseText);
        if (xhr.status === 200){
            
            if (json_data["code"]===0){
                window.alert(json_data["data"]);
                var ul = document.getElementById("item");
                var li = document.getElementById(fileid);
                ul.removeChild(li)
            }else{
                window.alert(json_data["data"]);
            }
        }else{
            window.alert(json_data["data"]);
        }
    }
    xhr.send(formData)
    
}

async function display() {

    var xhr = new XMLHttpRequest();
    xhr.open("GET", "/api/allfile", true);
    xhr.onload = function (){
        var json_data = JSON.parse(xhr.responseText);
        if (xhr.status === 200){

            if (json_data["code"]===0){
                updatepage(json_data["data"])
            }else{
                window.alert("错误：data"+json_data["data"]);
            }

        }else{
            window.alert("错误：data"+json_data["data"]);
        }
    }
    xhr.send()
}

window.onload = display()