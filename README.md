# reproducible-research-IA369Z

# A Predictive DASH QoE Approach Based on Machine Learning at Multi-access Edge Computing
# Abstract
Network operators expected to run on with fast growing multimedia video streaming (Dynamic Adaptive Streaming over HTTP, aka MPEG-DASH, such as YouTube and  Netflix) traffic demand while providing a high Quality of Experience (QoE) to the end-users. The cost, complexity, and scalability of existing QoE estimation solutions have significant limitations to infer QoE from network traffic. However, this works provides an end-user QoE estimation method based on a predictive passive QoE probe mechanism of DASH video using a Machine Learning (ML) approach running at network edge nodes. This work describes the design and implementation of probe configuration at the target edge, with traffic flow monitoring to generate network-level Quality of Service (QoS) metrics. Moreover, build a QoS-QoE correlation ML model in a real-time fashion to detect user equipment traffic patterns to predict user QoE more specifically Mean Opinion Score (MOS).

- [PDF version of paper](https://github.com/sajibtariq/reproducible-research-IA369Z/blob/master/deliver/project.pdf)
- [Executable paper](https://github.com/sajibtariq/reproducible-research-IA369Z/blob/master/deliver/third-draft.ipynb)

## Workflow
![alt text](https://github.com/sajibtariq/reproducible-research-IA369Z/blob/master/figures/Project%20workflow.jpg?raw=true)

## Requirements

**Data Acquisition**

[Testbed](https://github.com/sajibtariq/reproducible-research-IA369Z/tree/master/testbed)

**Data Pre-processing, Data Analysis & Executable Paper**

[Jupyter Notebook followed by anaconda and given instructions](https://github.com/sajibtariq/reproducible-research-IA369Z/blob/master/dev/anaconda_installation.txt)

**Note:** It is recommended install anaconda in home directory (e.g., /home/HOSTNAME/anaconda3)

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

* /data : Exported CSV along with splited training and testing CSV file
* /deliver : Executable notebook and original pdf
* /dev : Notebook along with data pre-process and analysis codes
* /figures : Research related figures
* /testbed : Data generation testbed aligned with codes
* /tesbed/experiment : Raw network data & video log storage

## Cautions

- Though there is three option available at below to execute the work. I will recommend to follow [option 3](#option-3--execute-the-work) to get the same result I presented at my work

- Option 1 and 2 shared for being fully transparent about the work and data provenance.  While fully reproduce, followed by option 1  manually testbed setup, data generation, and raw data, pre-process steps require a lot of time and storage to finish. You might get any unwanted error during manually testbed setup.
 - To avoid the dependency hell for option 1, a pre-built VM provided in option 2.  Testbed installation and raw data generation require a bit more memory space. Therefore fully reproducible pre-built VM is too large. To download this VM, good internet speed, and later work with it, a high configuration computer system recommended. 
- In options 1 and 2, you will get new raw data which pattern might be different from my data as I present in my work. Thus you may get a different final result.


##  Option 1 : Execute the Work

**Fully reproducible in Local Machine - Ubuntu 18.04**

**Note:** Testbed setup, data generation, and raw data pre-process steps require a lot of time and storage to finish. If you want you could skip this option and follow the option 2 or 3.


* **Environment Setup**

Step 1: [Testbed Setup](https://github.com/sajibtariq/reproducible-research-IA369Z/tree/master/testbed)

Step 2: [Jupyter Notebook followed by anaconda and given instructions](https://github.com/sajibtariq/reproducible-research-IA369Z/blob/master/dev/anaconda_installation.txt)

Step 3: [Scapy](https://anaconda.org/conda-forge/scapy)

Step 4: [Imbalanced-learn- python package](https://anaconda.org/conda-forge/imbalanced-learn)

Step 5: [nest-asyncio](https://anaconda.org/conda-forge/nest-asyncio)


* **Data Generation**

Step  1: Open terminal and execute the following command:

```bash
 cd ~/reproducible-research-IA369Z/testbed/

 sudo python3 test_1.py
```


* **Data Preprocess** 

Step 1: Open jupyter notebook and open the **third-draft.ipynb** notebook from **~/reproducible-research-IA369Z/deliver/** directory.


Step 2: To export csv from raw data, run the following cell in notebook:


```bash

%run ./Raw-data-preprocess-&-csv-export.ipynb
```

Step 3: To split training and testing csv from exported csv, run the following cell in notebook:

```bash
%run ~/reproducible-research-IA369Z/dev/preprocess.ipynb
```


* **Data Analysis**


Step 1: To execute the grid serach and model accuracy program, run the following cell in notebook:

```bash
%run ~/reproducible-research-IA369Z/dev/model_accuracy_with_grid_search.ipynb
```

Step 2: To get the prediction result, run the following cell in notebook:

```bash
%run ~/reproducible-research-IA369Z/dev/rfc.ipynb
```
 


##  Option 2 : Execute the Work


**Fully/Partially reproducible by provided pre-built VM including all dependencies**

[Pre-built VM ](#Distribution) is available in .ova format and should be usable with any modern x64 virtualization system.


**Note:** If you can download the  VM image and import successfully. But later,  not able to perform Data Generation and Raw data preprocess steps inside VM properly due to any technical problem. Then, skip the Data Generation and Data Preprocess (Raw data to export CSV) steps. You can still partially execute the work by completing  only  Data Preprocess (Split training and testing CSV from exported CSV)  and Data analysis steps with given dataset in  **~/reproducible-research-IA369Z/data/** directory 


* **Data Generation**

Step  1: Open terminal and execute the following command:

```bash
 cd ~/reproducible-research-IA369Z/testbed/

 sudo python3 test_1.py
 
```
* **Data Preprocess** 

Step 1: Open jupyter notebook and open the **third-draft.ipynb** notebook from **~/reproducible-research-IA369Z/deliver/** directory.


Step 2: Rename the 'Raw data preprocess & csv export.ipynb' notebook from **~/reproducible-research-IA369Z/dev/** directory as 'Raw-data-preprocess-&-csv-export.ipynb'.

Step 3: To export csv from raw data, run the following cell in notebook:

```bash
%run ./Raw-data-preprocess-&-csv-export.ipynb
```
****  


Step 4: To split training and testing csv from exported csv, run the following cell in notebook:

```bash
%run ~/reproducible-research-IA369Z/dev/preprocess.ipynb
```
* **Data Analysis**


Step 1: To execute the grid serach and model accuracy program, run the following cell in notebook:

```bash
%run ~/reproducible-research-IA369Z/dev/model_accuracy_with_grid_search.ipynb
```

Step 2: To get the prediction result, run the following cell in notebook:

```bash
%run ~/reproducible-research-IA369Z/dev/rfc.ipynb
```
 

 


##  [Option 3 : Execute the Work](#Option-3-:-Execute-the-Work)

**Reprducible based on pre-processed data from raw data (more specifically, skipping testbed setup, data generation, and raw data preprocessing steps)**

### Local machine:

* **Envirionment Setup** ubuntu 18.04 recommended

Step 1: [Jupyter Notebook followed by anaconda and given instructions](https://github.com/sajibtariq/reproducible-research-IA369Z/blob/master/dev/anaconda_installation.txt)

**Note:** It is recommended install anaconda in home directory (e.g., /home/HOSTNAME/anaconda3)

```bash
- Notebook==6.0.3
- Numpy==1.18.1
- Pandas==1.0.1
- Matplotlib==3.1.3
- Seaborn==0.10.0
- Scikit-Learn==0.22.1
- Scipy==1.4.1
```


Step 2: [Imbalanced-learn- python package](https://anaconda.org/conda-forge/imbalanced-learn)


* **Clone repository** 

Step 1: Open terminal and execute the following command:

```bash
git clone https://github.com/sajibtariq/reproducible-research-IA369Z.git

```

* **Data Preprocess** 

Step 1: Open jupyter notebook and open the **third-draft.ipynb** notebook from **~/reproducible-research-IA369Z/deliver/** directory.

Step 2: To split training and testing csv from exported csv, run the following cell in notebook:

```bash
%run ~/reproducible-research-IA369Z/dev/preprocess.ipynb
```

* **Data Analysis**


Step 1: To execute the grid serach and model accuracy program, run the following cell in notebook:

```bash
%run ~/reproducible-research-IA369Z/dev/model_accuracy_with_grid_search.ipynb
```

Step 2: To get the prediction result, run the following cell in notebook:

```bash
%run ~/reproducible-research-IA369Z/dev/rfc.ipynb

```
### [Pre-built VM:](#Distribution)

Wrapped with the packeages that just required for working with pre-proccess exported csv file.

* **To Update the Existing Project git Repository** 

Step 1: Open terminal and execute the following command:

```bash

cd  ~/reproducible-research-IA369Z

git pull

```
* **Data Preprocess** 

Step 1: Open jupyter notebook and open the **third-draft.ipynb** notebook from **~/reproducible-research-IA369Z/deliver/** directory.

Step 2: To split training and testing csv from exported csv, run the following cell in notebook:

```bash
%run ~/reproducible-research-IA369Z/dev/preprocess.ipynb
```

* **Data Analysis**

Step 1: To execute the grid serach and model accuracy program, run the following cell in notebook:

```bash
%run ~/reproducible-research-IA369Z/dev/model_accuracy_with_grid_search.ipynb
```

Step 2: To get the prediction result, run the following cell in notebook:

```bash
%run ~/reproducible-research-IA369Z/dev/rfc.ipynb
```


## Distribution

- [x] [VM](https://drive.google.com/open?id=1lwCD_fe47DXEOuD1L1LbO8a6otn_8WTq)[40 GB Size]: Ubuntu 18.04 x64 - Dash (pass: dash)

- [ ] Docker: Todo (Future Task)

## License
[GPL-3.0](https://github.com/sajibtariq/reproducible-research-IA369Z/blob/master/LICENSE)
