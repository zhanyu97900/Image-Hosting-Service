document.getElementById('registerForm').addEventListener('submit', function(event) {
    event.preventDefault(); // 阻止默认的表单提交行为
    Register()
    
});
var username = document.getElementById("username");
var password = document.getElementById("password");
var confirmPassword = document.getElementById("confirmPassword");
confirmPassword.addEventListener('input', checkPasswordsMatch);

function checkPasswordsMatch() {
  // 检查密码是否匹配
  if(password.value !== confirmPassword.value) {
    message.style.color = 'red';
    message.textContent = '密码不匹配';
  } else {
    message.style.color = 'green';
    message.textContent = '密码匹配';
  }
}

function Register(){
    usernameValue = username.value;
    passwordValue =password.value;
    var passwordNew = sha256Encrypt(passwordValue);
    var formData = new FormData();
    formData.append("username",usernameValue);
    formData.append("password",passwordNew);
    var statusRegister = document.getElementById("status");
    var xhr = new XMLHttpRequest()
    xhr.open("POST","/api/register",true);
    xhr.onload = function (){
        if (xhr.status===200){
            try {
                var registerJson = JSON.parse(xhr.responseText);
                if (registerJson["code"]==0){
                    statusRegister.textContent=registerJson["data"];
                    window.location="./login";
                }else{
                    statusRegister.textContent=registerJson["data"];
                }
            } catch (parseError) {
                console.log(parseError)
                window.alert("错误")
            }

        }else if (xhr.status===500){
            try {
                var registerJson = JSON.parse(xhr.responseText);
                if (registerJson["code"]==1){
                    statusRegister.textContent="服务器错误";
                }else{
                    statusRegister.textContent="未知服务器错误";
                }
                console.log(registerJson["data"])
            } catch (parseError) {
                console.log(parseError)
                window.alert("错误")
            }
        }else if (xhr.status===400){
            try {
                var registerJson = JSON.parse(xhr.responseText);
                if (registerJson["code"]==1){
                    
                    statusRegister.textContent=registerJson["data"];
                }else{
                    statusRegister.textContent=registerJson["data"];
                }
            } catch (parseError) {
                console.log(parseError)
                window.alert("错误")
            }}
    }
    xhr.send(formData);

}