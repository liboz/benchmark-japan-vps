package main

import benchmarkv1 "liboz/benchmark-japan-vps/gen/service/v1"

func Ptr(v int64) *int64 {
	return &v
}

var SAMPLE_RESULT = &benchmarkv1.BenchmarkResult{
	StartTime: 1662842149, EndTime: 1662842750, IoSpeed: 876300000,
	SpeedTestResults: []*benchmarkv1.SpeedTestResult{
		{City: "Los Angeles", Country: "US", UploadSpeed: 171260000, DownloadSpeed: 403810000},
		{City: "Amsterdam", Country: "NL", UploadSpeed: 117970000, DownloadSpeed: 621610000},
		{City: "Shanghai", Country: "CN", UploadSpeed: 48950000, DownloadSpeed: 60000},
		{City: "Singapore", Country: "SG", UploadSpeed: 47160000, DownloadSpeed: 6960000},
		{City: "Tokyo", Country: "JP", UploadSpeed: 69660000, DownloadSpeed: 650260000}},
	SingleCoreGeekbench: 3420,
	MultiCoreGeekbench:  3360,
	PingTestResults: []*benchmarkv1.PingTestResult{
		{Url: "google.com", DroppedPackets: Ptr(0), MinimumPing: 7.482, AveragePing: 7.897, MaximumPing: 10.613, StandardDeviation: 0.467},
		{Url: "youtube.com", DroppedPackets: Ptr(0), MinimumPing: 7.641, AveragePing: 8.105, MaximumPing: 11.958, StandardDeviation: 0.785},
		{Url: "amazon.co.jp", DroppedPackets: Ptr(0), MinimumPing: 75.103, AveragePing: 75.303, MaximumPing: 76.302, StandardDeviation: 0.186},
		{Url: "yahoo.co.jp", DroppedPackets: Ptr(0), MinimumPing: 182.567, AveragePing: 182.751, MaximumPing: 184.411, StandardDeviation: 0.294}}}
