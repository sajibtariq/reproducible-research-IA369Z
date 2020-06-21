# reproducible-research-IA369Z

# A Predictive DASH QoE Approach Based on Machine Learning at Multi-access Edge Computing
# Abstract
Network operators expected to run on with fast growing multimedia video streaming (Dynamic Adaptive Streaming over HTTP, aka MPEG-DASH, such as YouTube and Netflix) traffic demand while providing a high Quality of Experience (QoE) to the end-users. The cost, complexity, and scalability of existing QoE estimation solutions have significant limitations to infer QoE from network traffic. However, this works provides an end-user QoE estimation method based on a predictive passive QoE probe mechanism of DASH video using a Machine Learning (ML) approach running at network edge nodes. This work describes the design and implementation of probe configuration at the target edge, with traffic flow monitoring to generate network-level Quality of Service (QoS) metrics. Moreover, build a QoS-QoE correlation ML model in a real-time fashion to detect user equipment traffic patterns to predict user QoE more specifically Mean Opinion Score (MOS).

- [PDF version of paper](https://github.com/sajibtariq/reproducible-research-IA369Z/blob/master/deliver/project.pdf)
- [Executable paper](https://github.com/sajibtariq/reproducible-research-IA369Z/blob/master/deliver/third-draft.ipynb)

## Workflow
![alt text](https://github.com/sajibtariq/reproducible-research-IA369Z/blob/master/figures/Project%20workflow.jpg?raw=true)

## Requirements
**Data Acquisition**

[Testbed](https://github.com/sajibtariq/reproducible-research-IA369Z/tree/master/testbed)

**Data Pre-processing, Data Analysis & Executable Paper**

[Jupyter Notebook followed by Anaconda](https://docs.anaconda.com/anaconda/install/)
```bash
- Notebook==6.0.3
- Numpy==1.18.1
- Pandas==1.0.1
- Matplotlib==3.1.3
- Seaborn==0.10.0
- Scikit-Learn==0.22.1
- Scipy==1.4.1
```

[Scapy](https://anaconda.org/conda-forge/scapy)

[Imbalanced-learn- python package](https://anaconda.org/conda-forge/imbalanced-learn)

[nest-asyncio](https://anaconda.org/conda-forge/nest-asyncio)


## Folder Structure Scheme
* /data : Exported CSV file during preprocess from raw data
* /deliver : Executable notebook 
* /dev : Notebook along with data pre-process and analysis codes
* /figures : Research related figures
* /testbed : Data generation testbed aligned with codes
* /tesbed/experiment : Raw network data & video log strorage

##  Option 1 : Execute the Work

**Fully Execute in local Machine - Ubuntu 18.04**

**Note:** Testbed setup, data generation, and raw data pre-process steps require a lot of time and storage to finish. If you want you could skip this option and follow the Option 2 or Option 3


* **Environment Setup**

Step 1: [Testbed Setup](https://github.com/sajibtariq/reproducible-research-IA369Z/tree/master/testbed)

Step 2: [Jupyter Notebook followed by Anaconda](https://docs.anaconda.com/anaconda/install/)

Step 3: [Scapy](https://anaconda.org/conda-forge/scapy)

Step 4: [Imbalanced-learn- python package](https://anaconda.org/conda-forge/imbalanced-learn)

Step 5: [nest-asyncio](https://anaconda.org/conda-forge/nest-asyncio)


* **Data Generation**

Open terminal and execute the following command

$ cd ~/reproducible-research-IA369Z/testbed/

$ sudo python3 test_1.py


* **Data Preprocess** 

Step 1: Open jupyter notebook and open the **third-draft.ipynb** document from **~/reproducible-research-IA369Z/deliver/** directory


Step 2: **Raw data to export csv**

Run the cell that contains  **%run ./Raw-data-preprocess-&-csv-export.ipynb**  


Step 3: **Split training and testing CSV from exported csv**

Run the cell that contains  **%run ~/reproducible-research-IA369Z/dev/preprocess.ipynb** 


* **Data Analysis**

Step 1:  Run the cell that contains  **%run ~/reproducible-research-IA369Z/dev/model_accuracy_with_grid_search.ipynb**

Step 2:  Run the cell that contains  **%run ~/reproducible-research-IA369Z/dev/rfc.ipynb** 


##  Option 2 : Execute the Work


**Fully/Partially Exucute by provided pre-built VM images including all setup**

The VM images are in .ova format and should be usable with any modern x64 virtualization system.

[52 GB Size](https://drive.google.com/open?id=1lwCD_fe47DXEOuD1L1LbO8a6otn_8WTq) - Ubuntu 18.04 x64 - Dash (pass: dash)

**Note:** If you can download the  VM image and import successfully. But later,  not able to perform Data Generation and raw data preprocess steps inside VM properly due to any technical problem. Then, skip the Data Generation and Data Preprocess (Raw data to export CSV) steps. You can still partially execute the work by completing  only  Data Preprocess (Split training and testing CSV from exported CSV)  and Data analysis steps with given dataset in  **~/reproducible-research-IA369Z/data/** directory 


* **Data Generation**

Open terminal and execute the following command

$ cd ~/reproducible-research-IA369Z/testbed/

$ sudo python3 test_1.py

* **Data Preprocess** 

Step 1: Open jupyter notebook and open the  **third-draft.ipynb** document from **~/reproducible-research-IA369Z/deliver/** directory

Step 2: **Raw data to export csv**

Run the cell that contains  **%run ./Raw-data-preprocess-&-csv-export.ipynb**  


Step 3: **Split training and testing CSV from exported csv**


Run the cell that contains  **%run ~/reproducible-research-IA369Z/dev/preprocess.ipynb**  


* **Data Analysis**

Step 1:  Run the cell that contains  **%run ~/reproducible-research-IA369Z/dev/model_accuracy_with_grid_search.ipynb**

Step 2:  Run the cell that contains  **%run ~/reproducible-research-IA369Z/dev/rfc.ipynb** 


##  Option 3 : Execute the Work

**Partially executable if not able to perform Option 1/2 or want to skip Option 1/2 (more specifically, skipping testbed setup, data generation, and raw data preprocessing steps)**


* **Envirionment Setup** ubuntu 18.04 recommended

Step 1: [Jupyter Notebook followed by Anaconda](https://docs.anaconda.com/anaconda/install/)

Step 2: [Scappy](https://anaconda.org/conda-forge/scapy)

Step 3: [Imbalanced-learn- python package](https://anaconda.org/conda-forge/imbalanced-learn)

Step 4: [nest-asyncio](https://anaconda.org/conda-forge/nest-asyncio)



* **Clone repository** 

Open terminal and execute the following command

$ git clone https://github.com/sajibtariq/reproducible-research-IA369Z.git


* **Data Preprocess**

Step 1: Open jupyter notebook and open the **third-draft.ipynb** document from **~/reproducible-research-IA369Z/deliver/** directory


Step 2: **Split training and testing CSV from exported csv**

Run the cell that contains  **%run ~/reproducible-research-IA369Z/dev/preprocess.ipynb** 



* **Data Analysis**


Step 1:  Run the cell that contains  **%run ~/reproducible-research-IA369Z/dev/model_accuracy_with_grid_search.ipynb** 

Step 2:  Run the cell that contains  **%run ~/reproducible-research-IA369Z/dev/rfc.ipynb** 



## Distribution

- [VM](https://drive.google.com/open?id=1lwCD_fe47DXEOuD1L1LbO8a6otn_8WTq) :Ubuntu 18.04 x64 - Dash (pass: dash)
- Docker: Todo (Future Task)
