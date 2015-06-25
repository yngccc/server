(function () {
    var recent_posts = document.getElementById("dat-recent-posts");
    if (recent_posts) {
	if (recent_posts.lastElementChild) {
	    recent_posts.lastElementChild.style.borderBottomWidth = 0;
	}
    }
})();

(function () {
    var replys = document.getElementsByClassName("dat-reply");
    for (var i = 0; i < replys.length; i++) {
	replys[i].onclick = function(reply) {
	    return function() {
		var form = reply.nextElementSibling;
		form.hidden = !form.hidden;
	    }
	}(replys[i]);
    }
})();

(function () {
    var canvas = document.getElementById("dat-canvas");
    var ctx = canvas.getContext("2d");
    var canvas_div = document.getElementById("dat-canvas-div");
    canvas.width = canvas_div.offsetWidth;
    canvas.height = canvas_div.offsetHeight;

    window.setInterval(function() {
    	// dpendulum_rk4_step(1);
    	// console.log(dpendulum_w1, dpendulum_w2, dpendulum_a1, dpendulum_a2);
    	ctx.clearRect(0, 0, canvas.width, canvas.height);
    	ctx.beginPath();
    	ctx.arc(canvas.width / 2, canvas.height / 2, Math.random() * canvas.height / 2, 0, 2 * Math.PI);
    	ctx.stroke();
    }, 64);
})();

// // double pendulum
// var dpendulum_w1 = 0;
// var dpendulum_w2 = 0;
// var dpendulum_a1 = Math.PI / 2;
// var dpendulum_a2 = Math.PI / 2;
// var dpendulum_g_l = 1;

// var dpendulum_fw1 = function(w1, w2, a1, a2) {
// 	var numerator = w1 * w1 * Math.sin(2 * a1 - 2 * a2) + 2 * w2 * w2 * Math.sin(a1 - a2) + dpendulum_g_l * Math.sin(a1 - 2 * a2) + 3 * Math.sin(a1);
// 	var denominator = 3 - Math.cos(2 * a1 - 2 * a2);
// 	return -numerator / denominator;
// }

// var dpendulum_fw2 = function(w1, w2, a1, a2) {
// 	var numerator = 4 * w1 * w1 * Math.sin(a1 - a2) + w2 * w2 * Math.sin(2 * a1 - 2 * a2) + 2 * dpendulum_g_l * Math.sin(2 * a1 - a2) - Math.sin(a2);
// 	var denominator = 3 - Math.cos(2 * a1 - 2 * a2);
// 	return numerator / denominator;
// }

// var dpendulum_rk4_step = function(h) {
// 	var k1 = dpendulum_fw1(dpendulum_w1, dpendulum_w2, dpendulum_a1, dpendulum_a2);
// 	var l1 = dpendulum_fw2(dpendulum_w2, dpendulum_w2, dpendulum_a1, dpendulum_a2);

// 	var k2 = dpendulum_fw1(dpendulum_w1 + k1 * h / 2, dpendulum_w2 + l1 * h / 2, dpendulum_a1, dpendulum_a2);
// 	var l2 = dpendulum_fw2(dpendulum_w1 + k1 * h / 2, dpendulum_w2 + l1 * h / 2, dpendulum_a1, dpendulum_a2);

// 	var k3 = dpendulum_fw1(dpendulum_w1 + k2 * h / 2, dpendulum_w2 + l2 * h / 2, dpendulum_a1, dpendulum_a2);
// 	var l3 = dpendulum_fw2(dpendulum_w1 + k2 * h / 2, dpendulum_w2 + l2 * h / 2, dpendulum_a1, dpendulum_a2);

// 	var k4 = dpendulum_fw1(dpendulum_w1 + k3, dpendulum_w2 + l3, dpendulum_a1, dpendulum_a2);
// 	var l4 = dpendulum_fw2(dpendulum_w1 + k3, dpendulum_w2 + l3, dpendulum_a1, dpendulum_a2);

// 	dpendulum_w1 = dpendulum_w1 + h / 6 * (k1 + 2 * k2 + 2 * k3 + k4);
// 	dpendulum_w2 = dpendulum_w2 + h / 6 * (l1 + 2 * l2 + 2 * l3 + l4);

// 	dpendulum_a1 = dpendulum_a1 + dpendulum_w1 * h;
// 	dpendulum_a2 = dpendulum_a2 + dpendulum_w2 * h;
// }
