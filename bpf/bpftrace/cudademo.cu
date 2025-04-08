//#! nvcc -cudart=shared -o cudademo cudademo.cu && ./cudademo

#include <iostream>
#include <vector>
#include <cuda_runtime.h>
#include <device_launch_parameters.h>

// CUDA核函数，用于执行数组相加
__global__ void vectorAdd(const float* A, const float* B, float* C, int numElements)
{
    int i = blockDim.x * blockIdx.x + threadIdx.x;
    if (i < numElements)
    {
        C[i] = A[i] + B[i];
    }
}

// 检查CUDA运行时API调用是否成功
inline void checkCudaError(cudaError_t err, const char* msg)
{
    if (err != cudaSuccess)
    {
        std::cerr << "CUDA错误: " << msg << " (" << cudaGetErrorString(err) << ")" << std::endl;
        exit(EXIT_FAILURE);
    }
}

int main()
{
    // 打印CUDA设备信息
    cudaDeviceProp prop;
    checkCudaError(cudaGetDeviceProperties(&prop, 0), "获取设备属性失败");
    std::cout << "使用设备: " << prop.name << std::endl;
    std::cout << "CUDA计算能力: " << prop.major << "." << prop.minor << std::endl;

    // 数组元素数量
    const int numElements = 50000;
    std::cout << "数组相加示例，每个数组包含 " << numElements << " 个元素" << std::endl;

    // 计算需要的线程块和每块线程数
    const int threadsPerBlock = 256;
    const int blocksPerGrid = (numElements + threadsPerBlock - 1) / threadsPerBlock;
    std::cout << "CUDA配置: " << blocksPerGrid << " 个线程块 x " 
              << threadsPerBlock << " 个线程/块" << std::endl;

    // 分配主机内存（使用固定内存提高性能）
    float *h_A, *h_B, *h_C;
    checkCudaError(cudaMallocHost(&h_A, numElements * sizeof(float)), "分配固定主机内存A失败");
    checkCudaError(cudaMallocHost(&h_B, numElements * sizeof(float)), "分配固定主机内存B失败");
    checkCudaError(cudaMallocHost(&h_C, numElements * sizeof(float)), "分配固定主机内存C失败");

    // 初始化主机数组
    for (int i = 0; i < numElements; ++i)
    {
        h_A[i] = static_cast<float>(rand()) / RAND_MAX;
        h_B[i] = static_cast<float>(rand()) / RAND_MAX;
    }

    // 分配设备内存
    float *d_A, *d_B, *d_C;
    checkCudaError(cudaMalloc(&d_A, numElements * sizeof(float)), "分配设备内存A失败");
    checkCudaError(cudaMalloc(&d_B, numElements * sizeof(float)), "分配设备内存B失败");
    checkCudaError(cudaMalloc(&d_C, numElements * sizeof(float)), "分配设备内存C失败");

    // 将数据从主机复制到设备
    std::cout << "将数据从主机内存复制到设备内存..." << std::endl;
    checkCudaError(cudaMemcpy(d_A, h_A, numElements * sizeof(float), cudaMemcpyHostToDevice), 
                  "复制数据A到设备失败");
    checkCudaError(cudaMemcpy(d_B, h_B, numElements * sizeof(float), cudaMemcpyHostToDevice), 
                  "复制数据B到设备失败");

    // 启动CUDA核函数
    std::cout << "启动CUDA核函数..." << std::endl;
    vectorAdd<<<blocksPerGrid, threadsPerBlock>>>(d_A, d_B, d_C, numElements);
    checkCudaError(cudaGetLastError(), "内核启动失败");

    // 等待设备完成计算
    std::cout << "等待设备完成计算..." << std::endl;
    checkCudaError(cudaDeviceSynchronize(), "设备同步失败");

    // 将结果从设备复制回主机
    std::cout << "将结果从设备内存复制回主机内存..." << std::endl;
    checkCudaError(cudaMemcpy(h_C, d_C, numElements * sizeof(float), cudaMemcpyDeviceToHost), 
                  "复制结果到主机失败");

    // 验证结果
    std::cout << "验证前5个结果..." << std::endl;
    for (int i = 0; i < 5; ++i)
    {
        std::cout << h_A[i] << " + " << h_B[i] << " = " << h_C[i] << std::endl;
    }

    // 释放设备内存
    checkCudaError(cudaFree(d_A), "释放设备内存A失败");
    checkCudaError(cudaFree(d_B), "释放设备内存B失败");
    checkCudaError(cudaFree(d_C), "释放设备内存C失败");

    // 释放主机内存
    checkCudaError(cudaFreeHost(h_A), "释放固定主机内存A失败");
    checkCudaError(cudaFreeHost(h_B), "释放固定主机内存B失败");
    checkCudaError(cudaFreeHost(h_C), "释放固定主机内存C失败");

    std::cout << "完成!" << std::endl;
    return 0;
}