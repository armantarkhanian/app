window.onload = function() {
    const centrifuge = new Centrifuge('ws://localhost:8000/connection/websocket?format=json');
    function drawText(text) {
        const div = document.createElement('div');
        div.innerHTML = text + '<br>';
        document.body.appendChild(div);
    }
    centrifuge.on('connect', function(ctx){
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

    centrifuge.on('publish', function(ctx) {
        const channel = ctx.channel;
        const payload = JSON.stringify(ctx.data);
        alert('Publication from server-side channel ' + channel + ": " + payload);
    });


    const sub = centrifuge.subscribe("chat", function(ctx) {
        document.getElementsByTagName("title")[0].innerHTML = JSON.stringify(ctx.data);
        drawText(JSON.stringify(ctx.data));
    });

    const input = document.getElementById("input");
    input.addEventListener('keyup', function(e) {
        if (e.keyCode === 13) {
            sub.publish(this.value);
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
