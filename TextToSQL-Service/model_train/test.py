import torch
print(torch.cuda.is_available())  # Should return True if GPU is available
print(torch.cuda.get_device_name(0))  # Should print your GPU name (e.g., "GeForce GTX 1650")
