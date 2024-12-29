async function calculateSHA256(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = async (event) => {
            const arrayBuffer = event.target.result;
            const hashBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer);
            const hashArray = Array.from(new Uint8Array(hashBuffer));
            const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
            resolve(hashHex);
        };
        reader.onerror = reject;
        reader.readAsArrayBuffer(file);
    });
}
async function prepara(file) {
    const fileName = file.name;
    const fileHash = await calculateSHA256(file);
    const formData_prepare = new FormData();
    formData_prepare.append("filename",fileName);
    formData_prepare.append("sha256",fileHash);
    return new Promise((resolve, reject) => {
        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/api/prepare", true);
        xhr.onload = function () {
            if (xhr.status === 200) {
                try {
                    const prepareJson = JSON.parse(xhr.responseText);
                    if (prepareJson["code"] === 0) {
                        const fileId = prepareJson["data"];
                        resolve([fileId, null, false]);
                    } else if (prepareJson["code"] === 2) {
                        const fileId = prepareJson["data"];
                        const newFileId = prepareJson["newFileid"];
                        resolve([fileId, newFileId, true]);
                    } else {
                        throw new Error(`上传是失败: ${prepareJson["code"]}`);
                    }
                } catch (parseError) {
                    reject(new Error('解释失败'));
                }
            } else {
                reject(new Error(`服务器错误,上传是失败: ${xhr.status}`));
            }
        };
        xhr.onerror = function () {
            reject(new Error('网络问题'));
        };
        xhr.send(formData_prepare);
    });
}
async function upload(file,fileId) {
    //upload 过程
    const formData = new FormData();
    formData.append('file', file);
    formData.append("fileid",fileId);
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/upload", true);
    xhr.onload = function (){
        if (xhr.status === 200){
            const json_data = JSON.parse(xhr.responseText);
            if (json_data["code"] === 0){
                const fileid = json_data["data"]
                uploadcontainer=document.getElementById("upload-container");
                uploadresult=document.getElementById("upload-result");
                uploadcontainer.style.display = "none";
                uploadresult.style.display = "";
                img_url = "./file/" + fileid;
                var image=uploadresult.querySelector("img");
                image.src = img_url;
                var lableimg=uploadresult.querySelector("p");
                lableimg.textContent = window.location.origin+img_url.replace(".","");
                var button = uploadresult.querySelector("button");
                button.style.display = "none"
            }else if (json_data["code"] === 1){
                const fileid = json_data["data"]
                uploadcontainer=document.getElementById("upload-container");
                uploadresult=document.getElementById("upload-result");
                uploadcontainer.style.display = "none";
                uploadresult.style.display = "";
                var lableimg=uploadresult.querySelector("p");
                lableimg.textContent = fileid;
                var button = uploadresult.querySelector("button");
                button.style.display = "none"
            }

        }
    }
    xhr.send(formData);
    
}

async function uploadImage() {
    const fileInput = document.getElementById('image-upload');
    const file = fileInput.files[0];
    if (!file) return;
    // prepara 过程
    try {
        const [fileId, newFileId, havenew] = await prepara(file);
        if (havenew){
            uploadcontainer=document.getElementById("upload-container");
            uploadresult=document.getElementById("upload-result");
            uploadcontainer.style.display = "none";
            uploadresult.style.display = "";
            img_url = "./file/" + fileId;
            var image=uploadresult.querySelector("img");
            image.src = img_url;
            var lableimg=uploadresult.querySelector("p");
            lableimg.textContent= "数据存在该图片地址为:"+window.location.origin+img_url.replace(".","");
            var button = uploadresult.querySelector("button");
            button.style.textContent = "继续上传";
            button.onclick =() => upload(file, newFileId)
        }else{
            await upload(file,fileId)
        }
    } catch (error) {
        window.alert(error)
    }

}