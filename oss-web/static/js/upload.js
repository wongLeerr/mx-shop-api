accessid = "";
accesskey = "";
host = "";
policyBase64 = "";
signature = "";
callbackbody = "";
filename = "";
key = "";
expire = 0;
g_object_name = "";
g_object_name_type = "";
now = timestamp = Date.parse(new Date()) / 1000;

function send_request() {
  let xmlhttp = new XMLHttpRequest();

  if (xmlhttp != null) {
    // 同源请求，避免端口/域名写死导致拿不到 token
    serverUrl = window.location.origin + "/v1/oss/token";

    xmlhttp.open("GET", serverUrl, false);
    xmlhttp.send(null);
    console.log("🐶🐶responseText>>>", xmlhttp.responseText);
    return xmlhttp.responseText;
  } else {
    alert("Your browser does not support XMLHTTP.");
  }
}

function check_object_radio() {
  var tt = document.getElementsByName("myradio");
  for (var i = 0; i < tt.length; i++) {
    if (tt[i].checked) {
      g_object_name_type = tt[i].value;
      break;
    }
  }
}

function get_signature() {
  now = timestamp = Date.parse(new Date()) / 1000;
  if (expire < now + 3) {
    body = send_request();
    var obj = JSON.parse(body);
    host = obj["host"];
    policyBase64 = obj["policy"];
    accessid = obj["accessid"];
    signature = obj["signature"];
    expire = parseInt(obj["expire"]);
    callbackbody = obj["callback"];
    key = obj["dir"];
    return true;
  }
  return false;
}

function random_string(len) {
  len = len || 32;
  var chars = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678";
  var maxPos = chars.length;
  var pwd = "";
  for (i = 0; i < len; i++) {
    pwd += chars.charAt(Math.floor(Math.random() * maxPos));
  }
  return pwd;
}

function get_suffix(filename) {
  pos = filename.lastIndexOf(".");
  suffix = "";
  if (pos != -1) {
    suffix = filename.substring(pos);
  }
  return suffix;
}

function calculate_object_name(filename) {
  if (g_object_name_type == "local_name") {
    g_object_name += "${filename}";
  } else if (g_object_name_type == "random_name") {
    suffix = get_suffix(filename);
    g_object_name = key + random_string(10) + suffix;
  }
  return "";
}

function get_uploaded_object_name(filename) {
  if (g_object_name_type == "local_name") {
    tmp_name = g_object_name;
    tmp_name = tmp_name.replace("${filename}", filename);
    return tmp_name;
  } else if (g_object_name_type == "random_name") {
    return g_object_name;
  }
}

function set_upload_param(up, filename, ret) {
  if (ret == false) {
    ret = get_signature();
  }
  g_object_name = key;
  if (filename != "") {
    suffix = get_suffix(filename);
    calculate_object_name(filename);
  }
  new_multipart_params = {
    key: g_object_name,
    policy: policyBase64,
    OSSAccessKeyId: accessid,
    success_action_status: "200", //让服务端返回200,不然，默认会返回204
    callback: callbackbody,
    signature: signature,
  };

  up.setOption({
    url: host,
    multipart_params: new_multipart_params,
  });

  up.start();
}

var uploader = new plupload.Uploader({
  runtimes: "html5,flash,silverlight,html4",
  browse_button: "selectfiles",
  //multi_selection: false,
  container: document.getElementById("container"),
  flash_swf_url: "lib/plupload-2.1.2/js/Moxie.swf",
  silverlight_xap_url: "lib/plupload-2.1.2/js/Moxie.xap",
  url: "http://oss.aliyuncs.com",

  filters: {
    mime_types: [
      //只允许上传图片和zip文件
      // 补全常见图片格式（很多相机/微信默认是 jpeg/webp/avif）
      { title: "Image files", extensions: "jpg,jpeg,png,gif,bmp,webp,avif" },
      { title: "Zip files", extensions: "zip,rar" },
    ],
    max_file_size: "10mb", //最大只能上传10mb的文件
    prevent_duplicates: true, //不允许选取重复文件
  },

  init: {
    PostInit: function () {
      document.getElementById("ossfile").innerHTML = "";
      document.getElementById("postfiles").onclick = function () {
        set_upload_param(uploader, "", false);
        return false;
      };
    },

    FilesAdded: function (up, files) {
      document
        .getElementById("console")
        .appendChild(document.createTextNode("\nFilesAdded: " + files.length));
      plupload.each(files, function (file) {
        document.getElementById("ossfile").innerHTML +=
          '<div id="' +
          file.id +
          '">' +
          file.name +
          " (" +
          plupload.formatSize(file.size) +
          ")<b></b>" +
          '<div class="progress"><div class="progress-bar" style="width: 0%"></div></div>' +
          "</div>";
      });
    },

    BeforeUpload: function (up, file) {
      check_object_radio();
      set_upload_param(up, file.name, true);
    },

    UploadProgress: function (up, file) {
      var d = document.getElementById(file.id);
      d.getElementsByTagName("b")[0].innerHTML =
        "<span>" + file.percent + "%</span>";
      var prog = d.getElementsByTagName("div")[0];
      var progBar = prog.getElementsByTagName("div")[0];
      progBar.style.width = 2 * file.percent + "px";
      progBar.setAttribute("aria-valuenow", file.percent);
    },

    FileUploaded: function (up, file, info) {
      if (info.status == 200) {
        document
          .getElementById(file.id)
          .getElementsByTagName("b")[0].innerHTML =
          "upload to oss success, object name:" +
          get_uploaded_object_name(file.name) +
          " 回调服务器返回的内容是:" +
          info.response;
      } else if (info.status == 203) {
        document
          .getElementById(file.id)
          .getElementsByTagName("b")[0].innerHTML =
          "上传到OSS成功，但是oss访问用户设置的上传回调服务器失败，失败原因是:" +
          info.response;
      } else {
        document
          .getElementById(file.id)
          .getElementsByTagName("b")[0].innerHTML = info.response;
      }
    },

    Error: function (up, err) {
      // 让错误更可见：很多时候是格式/大小不符合，用户以为“没反应”
      console.log("plupload error", err);
      if (err.code == -600) {
        document
          .getElementById("console")
          .appendChild(
            document.createTextNode(
              "\n选择的文件太大了,可以根据应用情况，在upload.js 设置一下上传的最大大小"
            )
          );
      } else if (err.code == -601) {
        document
          .getElementById("console")
          .appendChild(
            document.createTextNode(
              "\n选择的文件后缀不对,可以根据应用情况，在upload.js进行设置可允许的上传文件类型"
            )
          );
      } else if (err.code == -602) {
        document
          .getElementById("console")
          .appendChild(document.createTextNode("\n这个文件已经上传过一遍了"));
      } else {
        document
          .getElementById("console")
          .appendChild(
            document.createTextNode(
              "\nError: code=" +
                err.code +
                ", message=" +
                err.message +
                ", response=" +
                err.response
            )
          );
      }
    },
  },
});

uploader.init();
