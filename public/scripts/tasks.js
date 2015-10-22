
var ToDoList = React.createClass({
  getInitialState: function() {
    return {data: []};
  },
  componentDidMount : function() {
     this.loadCommentsFromServer();
    setInterval(this.loadCommentsFromServer, this.props.pollInterval);
  },
  loadCommentsFromServer: function() {
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      success: function(data) {
        this.setState({data: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  handleTaskSubmit: function(task) {
    var tasks = this.state.data;
    tasks.push(task);
    // 
    this.setState({data: tasks}, function() {
      $.ajax({
          url: this.props.url,
          dataType: 'json',
          type: 'POST',
          data: task,
          success: function(data) {
            this.setState({data: data});
          }.bind(this),
          error: function(xhr, status, err) {
            console.error(this.props.url, status, err.toString());
          }.bind(this)
        });
    });
  },
  handleTaskDelete: function(id) {
    console.log("handleTaskDelete " + id);

    var temp = this.state.data;
    delete temp[id]
    // console.log(temp)
    this.setState(temp)

    $.ajax({
          url: this.props.url,
          dataType: 'json',
          type: 'DELETE',
          data: { id : id},
          success: function(data) {
            this.setState({data: data});
          }.bind(this),
          error: function(xhr, status, err) {
            console.error(this.props.url, status, err.toString());
          }.bind(this)
        });
  }, 
  render: function() {
    return (
      <div id="toDoList">
        <p>checkboxes</p>
        <CheckboxList data11={this.bread} data={this.state.data}  onTaskDelete={this.handleTaskDelete}/>
        
        <p>footer</p>
        <AddItem   random="simple"  onTaskSubmit={this.handleTaskSubmit} />
      </div>
    );
  }
});


var CheckboxList = React.createClass({
  getInitialState: function() {
    return {data: []};
  },
  checkboxListHandleDelete: function(id) {
    this.props.onTaskDelete(id);
  },
  render: function() {
    var checkBoxes = this.props.data.map(function (e) {
      return (
          <Checkbox key={e.id} content={e.content} id={e.id}  checkboxListDeleteSubmit={this.checkboxListHandleDelete} />
        );
    }, this);
    var checkBoxesData = this.props.data;
    
    return (
      <div className="randomList">
        {checkBoxes}
      </div>
    );
  }
});


var Checkbox = React.createClass({
  handleDelete: function(e) {
    e.preventDefault();
    console.log("handleDelete");
    var id = this.props.id;
    this.props.checkboxListDeleteSubmit(id);
  },
  render: function() {
    return (
      <form>
        <input type="checkbox" onClick={this.handleDelete} />{this.props.content} <a href="" onClick={this.handleDelete}>X</a><br/>
      </form>
    );
  }
});


var AddItem = React.createClass({
  handleSubmit: function(e) {
    e.preventDefault();
    
    var content = this.refs.content.value.trim();
    var contentItem = {content : content, other : "nothing"};
    this.props.onTaskSubmit(contentItem);
    this.refs.content.value = '';

  },
  render: function() {
    return (
      <form className="AddItemForm" onSubmit={this.handleSubmit}>
        <input type="text" ref="content"></input>
        <input  type="submit" value="Post" />
      </form>
    );
  }
});



ReactDOM.render(
  <ToDoList   url="/api/tasks" pollInterval={10000} />,
  document.getElementById('content')
);