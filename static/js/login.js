
document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault(); // 阻止默认的表单提交行为

    Login()
});
function Login(){
    var username = document.getElementById("username").value;
    var  password = document.getElementById("password").value;
    var encryptpassword = sha256Encrypt(password);
    var status_login = document.getElementById("status"); 
    var formData = new FormData;
    formData.append("username",username);
    formData.append("password",encryptpassword);
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/login", true);
    xhr.onload = function (){
        if (xhr.status === 200){
            try {
                const loginJson = JSON.parse(xhr.responseText);
                if (loginJson["code"]===0){
                    status_login.textContent = loginJson["data"];
                    window.location="/";
                }else{
                   status_login.textContent = loginJson["data"];
                }

            } catch (parseError) {
                status_login.textContent = "错误";
                console.log(parseError)
            }
        }else if (xhr.status===500){
            const loginJson = JSON.parse(xhr.responseText);
            status_login.textContent = "服务器错误";
            console.log(loginJson["data"])
        }else if (xhr.status===400){
            const loginJson = JSON.parse(xhr.responseText);
            status_login.textContent = loginJson["data"];
            console.log(loginJson["data"])
        }else{
            status_login.textContent = "未知错误";
        }
    }
    xhr.send(formData);
}