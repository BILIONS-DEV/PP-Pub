const {useState} = React;

const LikeButton = () => {
    const [is_liked, setLiked] = useState(false)
    const [list,setList]=useState([
        { title: "sdfsdf", unique: 'abcdefgh' },
        { title: "asdfasd", unique: 'iklmnopq' },
        { title: "asdfasdfa", unique: 'rstuvxyz' },])

    let array = [
        { title: "thái đập trai", unique: 'abcdefgh' },
        { title: "thái bro", unique: 'iklmnopq' },
        { title: "ngon lành", unique: 'rstuvxyz' },
    ]
    const onCLick = () => {
        setLiked(!is_liked)
        setList(array)
    }
    return (
        <div className="test">
            <button className="btn btn-primary" onClick={onCLick}>Like</button>
            {is_liked ? (
                <div>Bạn đã like</div>
            ) : (<div>Ban chua like</div>)}
            {list && list.map((item,idx)=>{
                return(
                    <div key={idx}>
                        {item.title}
                    </div>
                )
            })}
        </div>
    )
}
const rootElement = document.querySelector('#root');
ReactDOM.render(<LikeButton/>, rootElement);