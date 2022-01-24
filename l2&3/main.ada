with Ada.Text_IO; use Ada.Text_IO;
with ada.numerics.discrete_random;

procedure lista3 is
  n : integer := 8;
  d : integer := 1;

  -- -- -- -- --
  -- Random
  type randRange is new Integer range 0..(n-1);
	package Rand_Int is new ada.numerics.discrete_random(randRange);
   	use Rand_Int;
   	gen : Generator;

  -- -- -- -- --
  -- Edges
  type Edge is record
    start : Integer;
    stop : Integer;
  end record;

  type edgesArray is array (0..(n-2+d)) of Edge;
  edgesA : edgesArray;

  procedure generateEdges is
  begin
    for i in 0..(n-2) loop
      edgesA(i) := (start => i, stop => i+1);
    end loop;
  end generateEdges;

  -- -- -- -- --
  -- Shortcuts
  procedure generateShortCuts is
    sztart : Integer;
    sztop : Integer;
    dd : Integer;
    existE : Integer;
    i : Integer;
  begin
    dd := 0; -- ilośc dodanych już krawedźi
    existE := 0;
    i := n - 1;
    while dd < d loop
      reset(gen);
      sztart := Integer(random(gen));
      sztop := Integer(random(gen));

      if sztart < sztop and then sztart /= 0 and then sztop /= n-1 then

        for j in 0..n-2+d loop
          if edgesA(j).start = sztart and then edgesA(j).stop = sztop then
            existE := 1;
          end if;
        end loop;

        if existE /= 1 then
          Put_Line(Integer'Image(sztart) & Integer'Image(sztop));
          edgesA(i) := (start => sztart, stop => sztop);
          i := i + 1;
          dd := dd + 1;
        end if;

        existE := 0;
      end if;
    end loop;
  end generateShortCuts;

  -- -- -- -- --
  -- Vertices
  type Ver is record
    id : Integer;
  end record;

  type nodeArray is array (0..(n-1)) of Ver;
  nodeA : nodeArray;

  procedure generateVertices is
  begin
    for i in 0..(n-1) loop
      nodeA(i) := (id => i);
    end loop;
  end generateVertices;

  -- -- -- -- --
  -- Nexts
  procedure Nexts (V : Integer) is
  begin
    for i in edgesA'Range loop
      if edgesA(i).start = V or edgesA(i).stop = V then
        Put_Line(Integer'Image(edgesA(i).start) & " " & Integer'Image(edgesA(i).stop));
      end if;
    end loop;
  end Nexts;

  type neighborhood is array (0..(n-1)) of Integer;

  function NextsVertices (V : Integer) return neighborhood is
    n : neighborhood;
    index : Integer;
  begin
    index := 0;

    -- "zerowanie" tablicy
    for i in n'Range loop
      n(i) := -1;
    end loop;


    for i in edgesA'Range loop
      if edgesA(i).start = V then
        n(index) := edgesA(i).stop;
        index := index + 1;
      end if;
      if edgesA(i).stop = V then
        n(index) := edgesA(i).start;
        index := index + 1;
      end if;
    end loop;

    -- for i in n'Range loop
    --   Put(Integer'Image(n(i)));
    -- end loop;
    return n;
  end NextsVertices;

  function numberOfNexts (V : Integer) return Integer is
    no : Integer;
  begin
    no := 0;
    for i in edgesA'Range loop
      if edgesA(i).start = V or edgesA(i).stop = V then
        no := no + 1;
      end if;
    end loop;
    return no;
  end numberOfNexts;

  type R is record
    nexthop: Integer;
    cost : Integer;
    changed : Boolean;
  end record;

  -- -- -- -- --
  -- Mutex
  protected type Mutex is
    entry Seize;
    procedure Release;
  private
    Owned : Boolean := False;
  end Mutex;

  protected body Mutex is
    entry Seize when not Owned is
    begin
      Owned := True;
    end Seize;
    procedure Release is
    begin
      Owned := False;
    end Release;
  end Mutex;

  -- -- -- -- --
  -- Abs
  function absolute(a : Integer; b : Integer) return Integer is
    begin
    if a - b > 0 then
      return (a-b);
    else
      return (b-a);
    end if;
  end absolute;

  type para is record
    from : integer;
    jpara : Integer;
    rijcostpara : Integer;
  end record;

  type riArray is array (0..(n-1)) of R;
  type rofr is array (0..(n-1)) of riArray;
  routingOfrouting : rofr;
  raportrouting : rofr;



  -- -- -- -- --
  -- Sender - declaration
  task type Sender is
    entry Run(id : Integer);
    entry finish;
  end Sender;

  senderTask : array (0..(n-1)) of Sender;

  -- -- -- -- --
  -- Receiver
  task type Receiver is
    entry Run(id : Integer);
    entry Input(p : para);
    entry finish;
  end Receiver;

  receiverTask : array (0..(n-1)) of Receiver;


  -- -- -- -- --
  -- Server printer
  task printer is
    entry print(msg : String);
  end printer;

  task body printer is
  begin
    loop
      select
        accept print(msg : String) do
          Put_Line(msg);
        end print;
      or
        delay Duration(4.0);
        Put_Line("---------PROGRAM POWINNIEN SIĘ ZAKOŃCZYĆ ------ DRUKWOANIE RAPORTÓW -----------------");

        for i in raportrouting'Range loop
          for j in raportrouting(i)'Range loop
            Put_Line("Z wierzchołka " & integer'image(i) & " do " & integer'image(j) & " jest " & integer'image(raportrouting(i)(j).cost) & " skoków");
          end loop;
        end loop;


        for i in senderTask'Range loop
          -- Put_Line(integer'image(i) & " finish");
          senderTask(i).finish;
          -- Put_Line(integer'image(i) & " finish");
          receiverTask(i).finish;
        end loop;
        Put_Line("??????????????????????? CZEMU TA LINIJKA NIE ZACHODZI ;( ?????????????????????????????????");
        exit;
      end select;
    end loop;
  end printer;


  task body Receiver is
    pinput : para;
    newcost : Integer;
    tmp : R;
    M : Mutex;
    idreceiver : Integer;
  begin
    loop
      select
        accept Run(id : Integer) do
          idreceiver := id;
          -- Put_Line("RUN Receiver");

          M.Seize;
          -- Put_Line(Integer'Image(idreceiver) & " moje startowe routing tables");
          -- for re in routingOfrouting(idreceiver)'Range loop
          --   Put_Line(Integer'Image(re) & " "  & Integer'Image(routingOfrouting(idreceiver)(re).nexthop) & " " &
          --     Integer'Image(routingOfrouting(idreceiver)(re).cost)
          --     & " "  & Boolean'Image(routingOfrouting(idreceiver)(re).changed));
          -- end loop;

          raportrouting(idreceiver) := routingOfrouting(idreceiver);
          M.Release;
        end Run;
      or
        accept Input(p : para) do
          pinput := p;
          newcost := pinput.rijcostpara + 1;

          Put_Line(Integer'image(idreceiver) & " otrzymałem pakiet od wierzchołka" &
            Integer'Image(pinput.from) & " o wartości j = " & integer'image(pinput.jpara) & " i cost =" &
            integer'Image(pinput.rijcostpara)
            );

          M.Seize;
          if newcost < routingOfrouting(idreceiver)(pinput.jpara).cost then
            -- Put_Line("---------------------------------------------------------------------------------------------------");
            -- Put_Line(" stara wartość wartosc ->" & integer'image(routingOfrouting(idreceiver)(pinput.jpara).cost) & " nowa wartosc ->" & Integer'image(newcost));
            tmp := (nexthop => pinput.from, cost => newcost, changed => TRUE);
            routingOfrouting(idreceiver)(pinput.jpara) := tmp;
            -- Put_Line("                                 robie zmiane");

            -- printer.print(Integer'image(idreceiver) & " table R");
            -- for i in routingOfrouting(idreceiver)'range loop
            --   Put_Line("--->" & integer'image(idreceiver) & " " & integer'image(i) & " " & integer'Image(routingOfrouting(idreceiver)(i).cost) & " " & Boolean'image(routingOfrouting(idreceiver)(i).changed));
            -- end loop;

            -- ostatni zapis do raportów
            raportrouting(idreceiver) := routingOfrouting(idreceiver);
          end if;
          M.Release;
        end Input;
      or
        accept finish;
          Put_Line("Koniec receivera" & Integer'Image(idreceiver));
          exit;
      -- or
      --   terminate;
      end select;
    end loop;
  end Receiver;

  -- -- -- -- --
  -- Sender - body
  task body Sender is
    sasiady : neighborhood;
    idsender : Integer;
    tmp : R;
    p : para;
    r : randRange;
    rf : Float;
    M : Mutex;
  begin
    loop
      select
        accept Run(id : Integer) do
          idsender := id;
          sasiady := NextsVertices(id);
        end Run;

        loop
          for j in routingOfrouting(idsender)'Range loop
            if routingOfrouting(idsender)(j).changed = TRUE then
              for l in sasiady'Range loop
                M.Seize;
                p := (from => idsender, jpara => j, rijcostpara => routingOfrouting(idsender)(j).cost);
                tmp := (nexthop => routingOfrouting(idsender)(j).nexthop,
                        cost => routingOfrouting(idsender)(j).cost,
                        changed => FALSE);
                routingOfrouting(idsender)(j) := tmp;
                -- M.Release;

                for s in sasiady'Range loop
                  if sasiady(s) /= -1 then
                    receiverTask(sasiady(s)).Input(p);
                    Put_Line(Integer'Image(idsender) & " Wysyłam pakiet " &
                      Integer'Image(p.jpara) & " " & Integer'image(p.rijcostpara) &
                      " do sasiada" & Integer'Image(sasiady(s)));
                  end if;
                end loop;
                M.Release;

              end loop;
              -- M.Release;
            end if;
          end loop;

          -- random time to sleep
          reset(gen);
          r := random(gen);
          rf :=  Float (r) * 0.2;
          delay Duration (rf); -- sleep
        end loop;
      or
        accept finish;
          Put_Line("Koniec sendera" & Integer'Image(idsender));
          exit;
      -- or
      --   terminate;
      end select;
    end loop;
    Put_Line("Wyszedł");
  end Sender;

  -- -- -- -- --
  -- Node
  procedure Node (id: Integer) is
    -- riA : riArray;
    sasiady : neighborhood;
  begin
    -- Put_Line("My id:" &  Integer'Image(id));
    sasiady := NextsVertices(id);
    for j in nodeA'Range loop
      -- Put_Line("sonsiady" & Integer'Image(j));
      -- Put_Line(Integer'Image(numberOfNexts(j)));
      -- if id /= j then
        for next in 0..numberOfNexts(id) loop
          -- Put_Line(Integer'Image(next));
          -- if j = numberOfNexts(next) then -- TO JEST ŹLE PRZECIEŻ!!!!
          if j = sasiady(next) then
            -- Put_Line("standart ri");
            routingOfrouting(id)(j) := (nexthop => j, cost => 1, changed => TRUE);
          else
            -- Put_Line("trzeba wyliczać");
            if id < j then
              routingOfrouting(id)(j) := (nexthop => id + 1, cost => absolute(id, j), changed => TRUE);
            else
              routingOfrouting(id)(j) := (nexthop => id - 1, cost => absolute(id, j), changed => TRUE);
            end if;
          end if;
        end loop;
      -- end if;
    end loop;

    -- for j in riA'Range loop
    --   Put_Line(Integer'Image(id) & Integer'Image(riA(j).nexthop) & Integer'Image(riA(j).cost));
    -- end loop;
    -- routingOfrouting(id) := riA;
    receiverTask(id).Run(id);
    senderTask(id).Run(id);
  end Node;



-- test : neighborhood;
begin
  -- generate all edges with shortcuts and vertices
  generateEdges;
  generateShortCuts;
  generateVertices;

  -- printing graphs
  Put_Line("EDGES");
  for i in edgesA'Range loop
    Put_Line(Integer'Image(edgesA(i).start) & " ->" & Integer'Image(edgesA(i).stop));
  end loop;

  Put_Line("VERTICES");
  for i in nodeA'Range loop
    -- Put_Line(Integer'Image(nodeA(i).id));
    -- Nexts(i);
    Node(i);
  end loop;

  -- for j in nodeA'Range loop
  --   test := NextsVertices(j);
  --   Put_Line("Sąsiady dla " & Integer'image(j));
  --   for i in test'Range loop
  --     if test(i) /= -1 then
  --       Put(Integer'Image(test(i)) & " ");
  --     end if;
  --   end loop;
  --   Put_Line("");
  -- end loop;

end lista3;
