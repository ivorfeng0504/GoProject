
//加密
function do_encrypt(str) {
    var before = new Date();
    var rsa = new RSAKey();
    var key = "010001";
    var public_key = "96483cb253ae62ffb8bbc3cd5f8fbf4bd3d51ebb32c992bd7649a371bf83ed649cdd3d18f4ea7405438d7393c0c38c0daf5032df5d744e6c189be918fbf937d261e78ef807e173c41dedb7bafb4c72fe00cbc7a677e87e8d972512810a897fd31c8dde1f6607d708ed3e764d35a85a51767d4005ec6935e9c597b397ef46844b";
    var res;
    rsa.setPublic(public_key, key);
    res = rsa.encrypt(str);
    var after = new Date();
    if (res) {
        return linebrk(hex2b64(res), 64);
    }
    else {
        return "";
    }
}