import pandas as pd
import matplotlib.pyplot as plt
from sklearn.datasets import make_blobs, make_moons, make_circles
import os

n_samples = 1000
random_state = 42
output_file = "clustering_datasets.xlsx"
plot_file = "data_visualization.png"

datasets = {
    "Blobs": make_blobs(n_samples=n_samples, centers=3, random_state=random_state),
    "Moons": make_moons(n_samples=n_samples, noise=0.05, random_state=random_state),
    "Circles": make_circles(n_samples=n_samples, noise=0.05, factor=0.5, random_state=random_state)
}

plt.figure(figsize=(15, 5))

for i, (name, (data, labels)) in enumerate(datasets.items(), 1):
    df = pd.DataFrame(data, columns=["X", "Y"])
    df["Label"] = labels
    
    with pd.ExcelWriter(output_file, engine='openpyxl', mode='a' if os.path.exists(output_file) else 'w') as writer:
        df.to_excel(writer, sheet_name=name, index=False, header=False)
    
    plt.subplot(1, 3, i)
    plt.scatter(data[:, 0], data[:, 1], c=labels, cmap='viridis', s=10)
    plt.title(name)
    plt.xlabel("X")
    plt.ylabel("Y")

plt.tight_layout()
plt.savefig(plot_file)
plt.close()

print(f"Data saved to {output_file}")
print(f"Visualization of data: {plot_file}")