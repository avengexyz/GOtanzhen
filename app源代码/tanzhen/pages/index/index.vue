<template>
	<view>
		<u-navbar :is-back="false" title="Go探针">
			<view class="navbar-right" slot="right" @click="logout">
				注销
			</view>
		</u-navbar>
		<view class="content">
			<!-- 如果数组为空 -->
			<view v-if="info==''">
				<u-card :title="title">
					<view class="" slot="body">
						<view class="u-body-item u-flex u-border-bottom u-col-between u-p-t-0">
							<u-circle-progress active-color="#2979ff" :percent="80">
								<view class="u-progress-content">
									<view class="u-progress-dot"></view>
									<text class='u-progress-info'>CPU</text>
								</view>
							</u-circle-progress>
							<u-circle-progress active-color="#2979ff" :percent="80">
								<view class="u-progress-content">
									<view class="u-progress-dot"></view>
									<text class='u-progress-info'>内存</text>
								</view>
							</u-circle-progress>
							<u-circle-progress active-color="#2979ff" :percent="80">
								<view class="u-progress-content">
									<view class="u-progress-dot"></view>
									<text class='u-progress-info'>硬盘</text>
								</view>
							</u-circle-progress>
						</view>
					</view>
				</u-card>
			</view>
			<!-- 如果数组不为空 -->
			<view v-else>
				<view v-for="i in info">
					<u-card :title="i.data.服务器名称">
						<view class="" slot="body">
							<view class="u-body-item u-flex u-border-bottom u-col-between u-p-t-0">
								<u-circle-progress active-color="#2979ff" :percent="i.data.cpu百分比">
									<view class="u-progress-content">
										<view class="u-progress-dot"></view>
										<text class='u-progress-info'>CPU</text>
									</view>
								</u-circle-progress>
								<u-circle-progress active-color="#2979ff" :percent="i.data.内存百分比">
									<view class="u-progress-content">
										<view class="u-progress-dot"></view>
										<text class='u-progress-info'>内存</text>
									</view>
								</u-circle-progress>
								<u-circle-progress active-color="#2979ff" :percent="i.data.硬盘百分比">
									<view class="u-progress-content">
										<view class="u-progress-dot"></view>
										<text class='u-progress-info'>硬盘</text>
									</view>
								</u-circle-progress>
							</view>
						</view>
					</u-card>
				</view>
			</view>
		</view>
	</view>

</template>

<script>
	export default {
		data() {
			return {
				info: [],
				title: "服务器数据加载中"
			}
		},
		onLoad() {
			this.getlist()

		},
		methods: {
			//获取列表方法
			getlist() {
				//获取缓存内容
				var that = this;
				var post = uni.getStorageSync('post')
				var key = uni.getStorageSync('key')
				setTimeout(function() {
					// post请求
					uni.request({
						//网址
						url: 'http://' + post + '/ajax',
						//请求头 一定要加 要不然请求的就不是post 而是option
						header: {
							'content-type': 'application/x-www-form-urlencoded',
							// 'content-type': 'application/json; charset=utf-8',
						},
						//请求参数
						data: {},
						//请求方法
						method: "GET",
						success: function(res) {
							that.info = res.data;
						}
					})
				}, 3000)
			},
			//注销方法
			logout() {
				//清空本地缓存
				uni.clearStorageSync();
				//跳转至内容页
				uni.navigateTo({
					url: '/pages/login/login',
				})
			}
		}
	}
</script>

<style>
	page {
		background-color: #F6F6F6;
	}

	.u-card-wrap {
		background-color: $u-bg-color;
		padding: 1px;
	}

	.u-body-item {
		font-size: 32rpx;
		color: #333;
		padding: 20rpx 10rpx;
	}

	.u-progress-content {
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.u-progress-dot {
		width: 16rpx;
		height: 16rpx;
		border-radius: 50%;
		background-color: #fb9126;
	}

	.u-progress-info {
		font-size: 28rpx;
		padding-left: 16rpx;
		letter-spacing: 2rpx
	}

	.navbar-right {
		margin-right: 24rpx;
		display: flex;
	}
</style>
