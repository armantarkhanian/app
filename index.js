window.onload = function() {
    const centrifuge = new Centrifuge('ws://127.0.0.1/connection/websocket?format=protobuf');    
    const encoder = new TextEncoder("utf-8");    
    const decoder = new TextDecoder("utf-8");

    function drawText(text) {
        const div = document.createElement('div');
        div.innerHTML = text + '<br>';
        document.body.appendChild(div);
    }
    centrifuge.on('connect', function(ctx) {
        drawText('Connected over ' + ctx.transport);
        centrifuge.rpc({"method": "like", "userID": "15", "photoID": "34"}).then(function(res) {
            console.log('rpc result', res);
        }, function(err) {
            console.log('rpc error', err);
        });
    });

    centrifuge.on('disconnect', function(ctx) {
        drawText('Disconnected: ' + ctx.reason);
    });

    centrifuge.on('publish', function(ctx) { // server side publication to #user channel
        let channel = ctx.channel;
        let payload = decoder.decode(ctx.data);
        alert('Publication from server-side channel ' + channel + ": " + payload);
    });

    const sub = centrifuge.subscribe("chat", function(ctx) {  
        data = decoder.decode(ctx.data);
        document.getElementsByTagName("title")[0].innerHTML = data;
        drawText(data);
    });

    const input = document.getElementById("input");
    input.addEventListener('keyup', function(e) {
        if (e.keyCode === 13) {
            value = this.value.trim();
            if (value == "") {
                return
            }
            const binaryData = encoder.encode(value);
            sub.publish(binaryData);
            input.value = '';
        }
    });
    // After setting event handlers – initiate actual connection with server.
    centrifuge.connect();
}


function setCookie(name,value,days) {
    var expires = "";
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days*24*60*60*1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + (value || "")+ expires + "; path=/";
}
function getCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for(var i=0;i < ca.length;i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1,c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
    }
    return "";
}
