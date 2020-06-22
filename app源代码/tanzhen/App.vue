<script>
	export default {
		onLaunch: function () {
			console.log('初始化应用')
			//获取缓存内容
			var post = uni.getStorageSync('post')
			var key = uni.getStorageSync('key')

			//如果缓存为空跳转至登录页
			if(post==""||key==""){
				//跳转至登录页
				uni.navigateTo({
					url: '/pages/login/login',
				})
			}
			// post请求
			uni.request({
				//网址
				url: 'http://' + post + '/check',
				//请求头 一定要加 要不然请求的就不是post 而是option
				header: {
					'content-type': 'application/x-www-form-urlencoded',
					// 'content-type': 'application/json; charset=utf-8',
				},
				//请求参数
				data: {
					key: key
				},
				//请求方法
				method: "POST",
				success: (res) => {

					if (res.data.status != 200) {
						//清空本地缓存
						uni.clearStorageSync();
						//跳转至登录页
						uni.navigateTo({
							url: '/pages/login/login',
						})
					}

					//跳转至内容页
					uni.navigateTo({
						url: '/pages/index/index',
					})
				}
			})
		}
	}
</script>

<style lang="scss">
	/*每个页面公共css */
	@import "uview-ui/index.scss";
</style>
