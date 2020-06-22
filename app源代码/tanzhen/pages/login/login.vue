<template>
	
	<view>
		<u-navbar :is-back="false" title="Go探针"></u-navbar>
		<view class="content">
			<u-cell-group>
				<u-field v-model="post" label="地址" placeholder="请填写您的服务器地址">
				</u-field>
				<u-field v-model="key" label="验证码" placeholder="请填写您的服务器验证码">
				</u-field>
			</u-cell-group>
			<u-button @click="click" type="primary">登录</u-button>
		</view>
	</view>

</template>

<script>
	export default {
		data() {
			return {
				post: '',
				key: ''
			}
		},
		methods: {
			click() {
				// post请求
				uni.request({
					//网址
					url: 'http://' + this.post + '/check',
					//请求头 一定要加 要不然请求的就不是post 而是option
					header: {
						'content-type': 'application/x-www-form-urlencoded',
						// 'content-type': 'application/json; charset=utf-8',
					},
					//请求参数
					data: {
						key: this.key
					},
					//请求方法
					method: "POST",
					success: (res) => {
						
						if (res.data.status != 200) {
							//showToast 显示提示窗api
							return uni.showToast({
								title: `登录失败`,
							})
						}
						//清空本地缓存
						uni.clearStorageSync();
						//将post及key写入缓存
						uni.setStorageSync('post', this.post)
						uni.setStorageSync('key', this.key)
						//跳转至内容页
						uni.navigateTo({
							url: '/pages/index/index',
						})
					}
				})
			}
		}
	}
</script>

<style>

</style>
