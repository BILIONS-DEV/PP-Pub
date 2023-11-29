const {useState, useEffect} = React;

const ReactFloor = () => {

    const [sizes, setSize] = useState([])
    // const [sizeSelected, setSizeSelected] = useState([])

    useEffect(() => {
        GetAllSize()
    }, []);

    const GetAllSize = () => {
        axios.get('/target/loadsize')
            .then(function (response) {
                // handle success
                setSize(response.data)
            })
            .catch(function (error) {
                console.log(error);
            })
    }

    const Selected = (id) => {
        let currentIdx = sizes.map(item => {
            return item.id
        }).indexOf(id)
        if (currentIdx !== -1) {
            let newSelected = [...sizes]
            newSelected[currentIdx].selected = true
            setSize(newSelected)
        }
    }

    const Delete = (id) => {
        let currentIdx = sizes.map(item => {
            return item.id
        }).indexOf(id)
        if (currentIdx !== -1) {
            let newSelected = [...sizes]
            newSelected[currentIdx].selected = false
            setSize(newSelected)
        }
    }

    const Clear = () => {
        let newSize = [...sizes]
        newSize.map(item => {
            item.selected = false
        })
        setSize(newSize)
    }
    return (
        <div id="TargetingBox" className="card rounded-0 shadow-dark-80 border border-gray-50 mb-5">
            <div className="d-flex align-items-center bg-dominant px-4 toggle-collapse">
                <h5 className="card-header-title font-weight-semibold m-2">TARGETING</h5>
                <div className="ms-auto pe-md-2 collapsed">
                    <div className="export-dropdown">
                        <a aria-controls="collapseTargetingBox" data-bs-target="#collapseTargetingBox"
                           data-bs-toggle="collapse" aria-expanded="true" className="btn btn-link px-3">
                            <svg className="ms-2" xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                                 viewBox="0 0 13 13">
                                <rect data-name="Icons/Tabler/Chevron Down background" width="16"
                                      height="16" fill="none"></rect>
                                <path
                                    d="M.214.212a.738.738,0,0,1,.952-.07l.082.07L7.1,5.989a.716.716,0,0,1,.071.94L7.1,7.011l-5.85,5.778a.738.738,0,0,1-1.034,0,.716.716,0,0,1-.071-.94l.071-.081L5.547,6.5.214,1.233A.716.716,0,0,1,.143.293Z"
                                    transform="translate(13 3.25) rotate(90)" fill="#1e1e1e"></path>
                            </svg>
                        </a>
                    </div>
                </div>
            </div>
            <div id="collapseTargetingBox"
                 className="collapse show"
                 data-bs-parent="#TargetingBox">
                <div className="card-body px-5 py-for-card">
                    <div className="row my-4">
                        <div className="d-flex">
                            <label className="col-form-label form-label form-label-lg w-input-group">
                                Target
                            </label>
                            <div className="col-sm-8">
                                <div className="accordion muze-collapes border border-bottom-0 border-gray-50"
                                     id="target">
                                    <div className="card px-4">
                                        <div className="card-header p-0" id="AdSizes">
                                            <h5 className="mb-0 d-grid">
                                                <button
                                                    className="btn btn-light btn-block text-start p-3 rounded-0 fs-6 collapsed"
                                                    type="button" data-bs-toggle="collapse"
                                                    data-bs-target="#collapseAdSizes"
                                                    aria-expanded="false"
                                                    aria-controls="collapseAdSizes">
                                                    <div className="row w-100">
                                                        <div className="col-6">Size</div>
                                                        <div className="col-auto align-items-center d-flex">
                                                            <svg xmlns="http://www.w3.org/2000/svg"
                                                                 width="16"
                                                                 height="16" fill="currentColor"
                                                                 className="bi bi-check-lg text-blue-600 me-2"
                                                                 viewBox="0 0 16 16">
                                                                <path
                                                                    d="M13.485 1.431a1.473 1.473 0 0 1 2.104 2.062l-7.84 9.801a1.473 1.473 0 0 1-2.12.04L.431 8.138a1.473 1.473 0 0 1 2.084-2.083l4.111 4.112 6.82-8.69a.486.486 0 0 1 .04-.045z"/>
                                                            </svg>
                                                            <span className="fs-6 fw-normal"
                                                                  style={{fontSize: "0.85rem !important"}}
                                                                  id="text_for_size">all ad sizes</span>
                                                        </div>
                                                    </div>
                                                </button>
                                            </h5>
                                        </div>
                                        <div id="collapseAdSizes" className="collapse" aria-labelledby="AdSizes"
                                             data-bs-parent="#target">
                                            <div className="card-body card-target lh-lg pt-3 px-0">
                                                <div className="row">
                                                    <div className="col-lg-6 col-xxl-6">
                                                        <div className="row mb-3 header-target">
                                                            <div className="col-md-12">
                                                                <div
                                                                    className="d-flex search-left border border-gray-300 rounded-2">
                                                                    <img src="/static/svg/icons/search@16.svg"
                                                                         alt="Search" className="ps-2"/>
                                                                    <input type="search"
                                                                           className="form-control border-0 pe-0 input-search-target"
                                                                           placeholder="Search size..."
                                                                           id="search_ad_size"/>
                                                                </div>
                                                            </div>
                                                        </div>
                                                        <div id="list_ad_size" className="box-target">
                                                            {sizes && sizes.map((item, idx) => {
                                                                if (!item.selected) {
                                                                    return (
                                                                        <div key={item.id}
                                                                             className="target-item border-bottom border-gray-50 ms-2 me-3"
                                                                             id={item.id}>
                                                                            <div
                                                                                className="list-group list-group-flush my-n3">
                                                                                <div className="list-group-item">
                                                                                    <div
                                                                                        className="d-flex flex-row align-items-center">
                                                                                        <div className="col p-0">
                                                                                            <span>{item.name}</span>
                                                                                        </div>
                                                                                        <div className="col-auto">
                                                                                            <button type="button"
                                                                                                    onClick={() => Selected(item.id)}
                                                                                                    className="btn d-flex align-items-center btn-outline-secondary btn-icon rounded-circle p-0 add_size">
                                                                                                <svg
                                                                                                    xmlns="http://www.w3.org/2000/svg"
                                                                                                    width="12"
                                                                                                    height="12"
                                                                                                    fill="currentColor"
                                                                                                    className="bi bi-plus-lg"
                                                                                                    viewBox="0 0 16 16">
                                                                                                    <path
                                                                                                        d="M8 0a1 1 0 0 1 1 1v6h6a1 1 0 1 1 0 2H9v6a1 1 0 1 1-2 0V9H1a1 1 0 0 1 0-2h6V1a1 1 0 0 1 1-1z"/>
                                                                                                </svg>
                                                                                            </button>
                                                                                        </div>
                                                                                    </div>
                                                                                </div>
                                                                            </div>
                                                                        </div>
                                                                    )
                                                                }
                                                            })}
                                                        </div>
                                                    </div>
                                                    <div className="col-lg-6 col-xxl-6 block_size">
                                                        <div className="row mb-3 header-target">
                                                            <div className="col-md-12">
                                                                <div className="d-flex">
                                                                    <button onClick={() => {
                                                                        Clear()
                                                                    }}
                                                                            className="btn btn-clear-target ms-auto remove_all_size">
                                                                        CLEAR
                                                                    </button>
                                                                </div>
                                                            </div>
                                                        </div>
                                                        <div className="box-target-selected box_ad_size">
                                                            {sizes && sizes.map((item, idx) => {
                                                                if (item.selected) {
                                                                    return (
                                                                        <div key={item.id}
                                                                             className="target-item border-bottom border-gray-50 ms-2 me-3 item_selected"
                                                                             id={item.id}>
                                                                            <div
                                                                                className="list-group list-group-flush my-n3">
                                                                                <div className="list-group-item">
                                                                                    <div
                                                                                        className="d-flex flex-row align-items-center">
                                                                                        <div className="col p-0">
                                                                                            <span>{item.name}</span>
                                                                                        </div>
                                                                                        <div className="col-auto">
                                                                                            <button type="button"
                                                                                                    onClick={() => Delete(item.id)}
                                                                                                    className="btn d-flex align-items-center btn-outline-danger btn-icon rounded-circle p-0">
                                                                                                <svg
                                                                                                    xmlns="http://www.w3.org/2000/svg"
                                                                                                    width="12"
                                                                                                    height="12"
                                                                                                    fill="currentColor"
                                                                                                    className="bi bi-dash-lg"
                                                                                                    viewBox="0 0 16 16">
                                                                                                    <path
                                                                                                        d="M0 8a1 1 0 0 1 1-1h14a1 1 0 1 1 0 2H1a1 1 0 0 1-1-1z"></path>
                                                                                                </svg>
                                                                                            </button>
                                                                                        </div>
                                                                                    </div>
                                                                                </div>
                                                                            </div>
                                                                        </div>
                                                                    )
                                                                }
                                                            })}
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}
const rootElement = document.querySelector('#floor_component');
ReactDOM.render(<ReactFloor/>, rootElement);